package models

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/k8s"
	"github.com/opensourceways/app-robot-server/logs"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InstanceOpt struct {
	PName    string `form:"pName" json:"p_name" binding:"required"`
	PConfig  string `form:"pConfig" json:"p_config"`
	PVersion string `form:"pVersion" json:"p_version" binding:"required"`
	Replicas int    `form:"replicas" json:"replicas" binding:"number,min=1,max=10"`
}

func (iOpt InstanceOpt) CreateInstanceRecord(uName string) global.Error {
	version, gErr := GetPluginVersionDetail(iOpt.PName, iOpt.PVersion)
	if gErr != nil {
		return gErr
	}
	instance := iOpt.transformDBInstance(uName, version)
	err := dbmodels.GetDB().AddInstance(instance)
	if err == nil {
		return nil
	}
	if err.IsErrorOf(dbmodels.ErrRecordExists) {
		return global.ResponseError{ErrCode: global.PluginNameExistCode, Reason: global.PluginNameExistMsg}
	}
	return global.NewResponseSystemError()
}

func (iOpt InstanceOpt) transformDBInstance(uName string, version dbmodels.PluginVersion) dbmodels.Instance {
	now := time.Now().Unix()
	ins := dbmodels.Instance{
		ID:        uuid.NewString(),
		PName:     iOpt.PName,
		PConfig:   iOpt.PConfig,
		PVersion:  iOpt.PVersion,
		CreateBy:  uName,
		CreateAt:  now,
		DReplicas: int32(iOpt.Replicas),
	}
	ins.Status = global.InsSavedNotRun
	ins.DName = getDeployName(iOpt.PName)
	ins.DLabel = iOpt.PName
	ins.DPort = version.Port
	if iOpt.Replicas <= 0 {
		ins.DReplicas = 1
	}
	ins.SName = getServiceName(iOpt.PName)
	ins.SPort = version.Port
	ins.SType = "ClusterIP"
	return ins
}

func GetInstanceDetail(insID string) (dbmodels.Instance, global.Error) {
	detail, idbError := dbmodels.GetDB().GetInstance(insID)
	if idbError != nil {
		logs.Logger.Error(idbError)
		if idbError.IsErrorOf(dbmodels.ErrNoDBRecord) {
			return detail, global.ResponseError{ErrCode: global.NoRecodeCode, Reason: global.NoRecodeMsg}
		}
		return detail, global.NewResponseSystemError()
	}

	return detail, nil
}

func StartPluginInstance(insID string) global.Error {
	detail, gErr := GetInstanceDetail(insID)
	if gErr != nil {
		return gErr
	}
	version, gErr := GetPluginVersionDetail(detail.PName, detail.PVersion)
	if gErr != nil {
		return gErr
	}
	// create service
	service, gErr := createService(&detail)
	if gErr != nil {
		return gErr
	}
	logs.Logger.Info(fmt.Printf("create service success %q. \n", service.GetObjectMeta().GetName()))
	// create deployment
	cd, gErr := createDeployment(&detail, &version)
	if gErr != nil {
		//delete the success service pods
		if service != nil {
			if err := deleteService(service); err != nil {
				logs.Logger.Error(err)
			}
		}
		return gErr
	}
	logs.Logger.Info(fmt.Printf("create deployment success %q. \n", cd.GetObjectMeta().GetName()))

	idbErr := dbmodels.GetDB().UpdateInstanceStatus(insID, global.InsRunning)
	if idbErr != nil {
		logs.Logger.Error(idbErr)
	}
	return nil
}

func deleteService(service *corev1.Service) error {
	return k8s.GetK8sClient().DeleteService(service.Name)
}

func DeletePluginInstance(insID string) global.Error {
	detail, gErr := GetInstanceDetail(insID)
	if gErr != nil {
		return gErr
	}
	if err := k8s.GetK8sClient().DeleteService(detail.SName); err != nil {
		logs.Logger.Error(err)
		if !errors.IsNotFound(err) {
			return global.NewResponseSystemError()
		}
	}

	if err := k8s.GetK8sClient().DeleteDeployment(detail.DName); err != nil {
		logs.Logger.Error(err)
		if !errors.IsNotFound(err) {
			return global.NewResponseSystemError()
		}
	}
	idbErr := dbmodels.GetDB().UpdateInstanceStatus(insID, global.InsSavedNotRun)
	if idbErr != nil {
		logs.Logger.Error(idbErr)
	}

	return nil
}

func genService(ins *dbmodels.Instance) *corev1.Service {
	if ins == nil {
		return nil
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: ins.SName,
		},
		Spec: corev1.ServiceSpec{
			Selector: getDLabelKey(ins.DLabel),
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     transformPort(ins.SPort),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: transformPort(ins.DPort),
					},
				},
			},
		},
	}
}

func genDeployment(ins *dbmodels.Instance, version *dbmodels.PluginVersion) *appsv1.Deployment {
	if ins == nil || version == nil {
		return nil
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   ins.DName,
			Labels: getDLabelKey(ins.DLabel),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(ins.DReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: getDLabelKey(ins.DLabel),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getDLabelKey(ins.DLabel),
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  ins.PName,
							Image: version.Image,
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: transformPort(ins.DPort),
								},
							},
						},
					},
				},
			},
		},
	}
}

func createService(ins *dbmodels.Instance) (*corev1.Service, global.Error) {
	service := genService(ins)
	if service == nil {
		return nil, global.NewResponseSystemError()
	}
	svc, err := k8s.GetK8sClient().CreateService(service)
	if err != nil {
		logs.Logger.Error(err)
		return nil, global.ResponseError{ErrCode: global.SystemErrorCode, Reason: err.Error()}
	}
	return svc, nil
}

func createDeployment(ins *dbmodels.Instance, version *dbmodels.PluginVersion) (*appsv1.Deployment, global.Error) {
	deployment := genDeployment(ins, version)
	if deployment == nil {
		return nil, global.NewResponseSystemError()
	}
	cd, err := k8s.GetK8sClient().CreateDeployment(deployment)
	if err != nil {
		logs.Logger.Error(err)
		return nil, global.ResponseError{ErrCode: global.SystemErrorCode, Reason: err.Error()}
	}
	return cd, nil
}

func getDeployName(pName string) string {
	return fmt.Sprintf("deploy-%s", pName)
}

func getServiceName(pName string) string {
	return fmt.Sprintf("svc-%s", pName)
}

func getDLabelKey(label string) map[string]string {
	return map[string]string{"app": label}
}

func int32Ptr(i int32) *int32 { return &i }

func transformPort(p string) int32 {
	if p == "" {
		return 80
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		logs.Logger.Error(err)
		return 80
	}
	return int32(port)
}
