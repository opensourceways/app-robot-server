package k8s

import (
	"context"
	"fmt"
	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	tav1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	tcv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/opensourceways/app-robot-server/config"
)

var k8sClient Client

func GetK8sClient() Client {
	return k8sClient
}

type Client interface {
	CreateNamespaceIfNotExist(ns string) error
	CreateDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error)
	GetDeployment(name string) (*appsv1.Deployment, error)
	ReplaceDeployment(deploy *appsv1.Deployment) error
	UpdateDeployment(name, image string, replicas int32) error
	DeleteDeployment(name string) error

	CreateService(svc *corev1.Service) (*corev1.Service, error)
	DeleteService(name string) error
}

type client struct {
	*kubernetes.Clientset
	ctx       context.Context
	namespace string
}

func (c *client) namespaceClient() tcv1.NamespaceInterface {
	return c.CoreV1().Namespaces()
}

func (c *client) deploymentsClient() tav1.DeploymentInterface {
	return c.AppsV1().Deployments(c.namespace)
}

func (c client) servicesClient() tcv1.ServiceInterface {
	return c.CoreV1().Services(c.namespace)
}

func (c *client) CreateNamespaceIfNotExist(ns string) error {
	_, err := c.namespaceClient().Get(c.ctx, ns, metav1.GetOptions{})
	if err == nil {
		return err
	}
	if !errors.IsNotFound(err) {
		return err
	}
	nso := corev1.Namespace{}
	nso.Name = ns
	nso.APIVersion = "v1"
	nso.Kind = "Namespace"
	_, err = c.namespaceClient().Create(c.ctx, &nso, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateDeployment(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	return c.deploymentsClient().Create(c.ctx, deploy, metav1.CreateOptions{})
}

func (c *client) GetDeployment(name string) (*appsv1.Deployment, error) {
	return c.deploymentsClient().Get(c.ctx, name, metav1.GetOptions{})
}

func (c *client) ReplaceDeployment(deploy *appsv1.Deployment) error {
	_, err := c.deploymentsClient().Update(c.ctx, deploy, metav1.UpdateOptions{})
	return err
}

func (c *client) UpdateDeployment(name, image string, replicas int32) error {
	if image == "" && replicas == 0 {
		return fmt.Errorf("can't update the %s deployment with empty image and zero replicas", name)
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, err := c.GetDeployment(name)
		if err != nil {
			return err
		}
		if replicas != 0 {
			deployment.Spec.Replicas = int32Ptr(replicas)
		}
		if image != "" {
			deployment.Spec.Template.Spec.Containers[0].Image = image
		}
		return c.ReplaceDeployment(deployment)
	})
}

func (c *client) DeleteDeployment(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return c.deploymentsClient().Delete(c.ctx, name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}

func (c *client) CreateService(svc *corev1.Service) (*corev1.Service, error) {
	return c.servicesClient().Create(c.ctx, svc, metav1.CreateOptions{})
}

func (c *client) DeleteService(name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return c.servicesClient().Delete(c.ctx, name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}

func Init(cfg *config.KubernetesConfig) error {
	kCfg, err := clientcmd.BuildConfigFromFlags("", cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	clientSet, err := kubernetes.NewForConfig(kCfg)
	if err != nil {
		return err
	}
	k8sClient = &client{Clientset: clientSet, ctx: context.TODO(), namespace: cfg.Namespace}

	return k8sClient.CreateNamespaceIfNotExist(cfg.Namespace)
}

func int32Ptr(i int32) *int32 { return &i }
