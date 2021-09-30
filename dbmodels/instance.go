package dbmodels

import "github.com/opensourceways/app-robot-server/global"

type Instance struct {
	ID        string `json:"id"`
	// Status represents the current instance status.
	// 0: Saved but not running
	// 1: already running
	// -1: start exception
	// 2: history running and need update pod
	Status    global.InstanceStatus `json:"status"`
	PName     string `json:"p_name" bson:"p_name"`
	PVersion  string `json:"p_version" bson:"p_version"`
	PConfig   string `json:"p_config" bson:"p_config"`
	DName     string `json:"d_name" bson:"d_name"`
	DLabel    string `json:"d_label" bson:"d_label"`
	DPort     string `json:"d_port" bson:"d_port"`
	// DReplicas deployment number of copies
	DReplicas int32   `json:"d_replicas" bson:"d_replicas"`
	SName     string `json:"s_name" bson:"s_name"`
	// SPort the port of service
	SPort     string `json:"s_port" bson:"s_port"`
	// SType the type of service ClusterIP or NodePort default ClusterIP
	SType     string `json:"s_type" bson:"s_type"`
	CreateAt  int64  `json:"create_at" bson:"create_at"`
	CreateBy  string `json:"create_by" bson:"create_by"`
	UpdateBy  string `json:"update_by" bson:"update_by"`
	UpdateAt  int64  `json:"update_at" bson:"update_at"`
}
