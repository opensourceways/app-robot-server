package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/middleware"
	"github.com/opensourceways/app-robot-server/models"
)

type InstanceController struct {
	baseController
}

func (ic *InstanceController) RegisterRouter(rg *gin.RouterGroup) {
	icRouter := rg.Group("/manage").Use(middleware.JWTAuth())
	{
		icRouter.POST("/instances", ic.InstallPlugin)
		icRouter.POST("/instances/:insID/options", ic.StartPluginInstance)
		icRouter.DELETE("/instances/:insID/options", ic.DeletePluginInstance)
	}
}

// @Tags	manage
// @Security ApiKeyAuth
// @Summary	create plugin metadata for k8s cluster
// @Produce	json
// @Accept x-www-form-urlencoded
// @Param	pName	formData 	string	true	"plugin name"
// @Param	pConfig formData	string	false	"plugin config"
// @Param	pVersion	formData	string true "plugin version"
// @param	replicas	formData	int	false	"Number of pod copies (1-10)"
// @Success	201	{object}	models.BaseResponse{}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	401 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /manage/instances [post]
func (ic *InstanceController) InstallPlugin(c *gin.Context) {
	var p models.InstanceOpt
	if err := c.ShouldBind(&p); err != nil {
		ic.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	userName := ic.GetUserName(c)
	if gErr := p.CreateInstanceRecord(userName); gErr != nil {
		ic.responseWithError(c, gErr)
	} else {
		ic.responseSuccessCreate(c, nil)
	}
}

func (ic *InstanceController) UnInstallPlugin(c *gin.Context) {

}

func (ic *InstanceController) UpdatePluginVersion(c *gin.Context) {

}

func (ic *InstanceController) UpdateInstanceConfig(c *gin.Context) {

}

// @Tags	manage
// @Security ApiKeyAuth
// @Summary	create instances pods on k8s cluster by plugin metadata
// @Produce	json
// @Accept x-www-form-urlencoded
// @Param	insID	path 	string	true	"plugin name"
// @Success	201	{object}	models.BaseResponse{}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	401 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /manage/instances/{insID}/options [post]
func (ic *InstanceController) StartPluginInstance(c *gin.Context) {
	insID := c.Param("insID")
	if insID == "" {
		ic.responseWithError(c, global.NewIllegalInputErr())
		return
	}
	if gErr := models.StartPluginInstance(insID); gErr != nil {
		ic.responseWithError(c, gErr)
	} else {
		ic.responseSuccessCreate(c, nil)
	}
}

// @Tags	manage
// @Security ApiKeyAuth
// @Summary	delete k8s cluster pods by insID
// @Produce	json
// @Accept x-www-form-urlencoded
// @Param	insID	path 	string	true	"plugin name"
// @Success	201	{object}	models.BaseResponse{}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	401 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /manage/instances/{insID}/options [delete]
func (ic *InstanceController) DeletePluginInstance(c *gin.Context) {
	insID := c.Param("insID")
	if insID == "" {
		ic.responseWithError(c, global.NewIllegalInputErr())
		return
	}

	if gErr := models.DeletePluginInstance(insID); gErr != nil {
		ic.responseWithError(c, gErr)
	} else {
		ic.responseSuccessCreate(c, nil)
	}

}
