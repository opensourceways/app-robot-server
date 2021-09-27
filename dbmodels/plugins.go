package dbmodels

import (
	"github.com/opensourceways/app-robot-server/global"
)

const (
	FieldPluginName   = "name"
	FieldPluginAuthor = "author"
	FieldPVNumber     = "version_number"
	FieldPVersions    = "versions"
)

type Plugin struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// Status plugin current status. 0: submitted 1: audited 2: published -1: unavailable
	Status        global.PluginStatus `json:"status,omitempty"`
	Author        string              `json:"author,omitempty" bson:"author_id"`
	Versions      []PluginVersion     `json:"versions,omitempty"`
	RepositoryURL string              `json:"repository_url,omitempty" bson:"repository_url"`
	LastVersion   string              `json:"last_version,omitempty" bson:"last_version"`
	Collaborators []string            `json:"collaborators,omitempty"`
	CreateTime    int64               `json:"create_time,omitempty" bson:"create_time"`
	UpdateTime    int64               `json:"update_time,omitempty" bson:"update_time"`
	AuditBy       string              `json:"audit_by,omitempty" bson:"audit_by"`
}

type PluginVersion struct {
	VersionNumber string `json:"version_number" bson:"version_number"`
	VersionID     int    `json:"version_id" bson:"version_id"`
	ConfigExample string `json:"config_example" bson:"config_example"`
	Image         string `json:"image"`
	Port          string `json:"port"`
	Description   string `json:"description"`
	UploadBy      string `json:"upload_by" bson:"upload_by"`
	UploadTime    int64  `json:"upload_time" bson:"upload_time"`
}
