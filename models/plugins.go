package models

import (
	"time"

	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/logs"
)

type Plugin struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description" binding:"required"`
	RepoURL       string   `json:"repo_url" binding:"required"`
	Collaborators []string `json:"collaborators"`
}

func (p Plugin) Save(userName string) global.Error {
	dbp := dbmodels.Plugin{
		Name:          p.Name,
		Description:   p.Description,
		RepositoryURL: p.RepoURL,
		Collaborators: p.Collaborators,
		Author:        userName,
		Status:        global.PluginStatusAudited,
		AuditBy:       userName,
		CreateTime:    time.Now().Unix(),
		Versions:      []dbmodels.PluginVersion{},
	}
	err := dbmodels.GetDB().AddPlugin(dbp)
	if err == nil {
		return nil
	}
	if err.IsErrorOf(dbmodels.ErrRecordExists) {
		return global.ResponseError{ErrCode: global.PluginNameExistCode, Reason: global.PluginNameExistMsg}
	}
	return global.NewResponseSystemError()
}

type PVersionOPT struct {
	//VersionNumber is current version number. e.g: v1.1.0
	VersionNumber string `json:"version_number" binding:"required"`
	//Image corresponding to the current version. e.g: mysql:1.1.1
	Image string `json:"image" binding:"required"`
	//Port that the program listens on
	Port string `json:"port" binding:"required"`
	//ConfigExample example of current program configuration file
	ConfigExample string `json:"config_example"`
	//Description of the changes in the current version
	Description string `json:"description"`
}

func (pv *PVersionOPT) AddVersion(pName, uName string) global.Error {
	pDetails, gErr := GetPluginDetails(uName, pName)
	if gErr != nil {
		return gErr
	}
	vId := len(pDetails.Versions) + 1
	dpv := dbmodels.PluginVersion{
		VersionNumber: pv.VersionNumber,
		VersionID:     vId,
		ConfigExample: pv.ConfigExample,
		Image:         pv.Image,
		Port:          pv.Port,
		Description:   pv.Description,
		UploadBy:      uName,
		UploadTime:    time.Now().Unix(),
	}
	idbErr := dbmodels.GetDB().AddPluginVersion(pName, dpv)
	if idbErr != nil {
		if idbErr.IsErrorOf(dbmodels.ErrNoDBRecord) {
			return global.ResponseError{ErrCode: global.PluginVersionIsExistCode, Reason: global.PluginVersionIsExistMsg}
		}
		return global.NewResponseSystemError()
	}
	// update the plugin last version and status
	needPublish := pDetails.Status == global.PluginStatusAudited
	if idbErr := dbmodels.GetDB().UpdatePluginLastVersion(pName, dpv.VersionNumber, needPublish); idbErr != nil {
		logs.Logger.Error(idbErr)
	}
	return nil
}

func GetPluginsByUser(userName string) ([]dbmodels.Plugin, global.Error) {
	plugins, idbError := dbmodels.GetDB().GetUserPlugins(userName)
	if idbError != nil && idbError.IsErrorOf(dbmodels.ErrNoDBRecord) {
		logs.Logger.Error(idbError)
		return nil, global.NewResponseSystemError()
	}

	return plugins, nil
}

func GetPluginDetails(userName, pluginName string) (dbmodels.Plugin, global.Error) {
	detail, idbError := dbmodels.GetDB().GetPluginDetail(pluginName, userName)
	if idbError == nil {
		return detail, nil
	}
	if idbError.IsErrorOf(dbmodels.ErrNoDBRecord) {
		return detail, global.ResponseError{ErrCode: global.NoRecodeCode, Reason: global.NoRecodeMsg}
	}
	logs.Logger.Error(idbError)
	return detail, global.NewResponseSystemError()
}
