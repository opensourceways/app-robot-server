package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/middleware"
	"github.com/opensourceways/app-robot-server/models"
)

type PluginsController struct {
	baseController
}

func (pc *PluginsController) RegisterRouter(rg *gin.RouterGroup) {
	pcRouter := rg.Group("/developer").Use(middleware.JWTAuth())
	{
		pcRouter.POST("/plugins", pc.Add)
		pcRouter.GET("/plugins", pc.GetUserPlugins)
		pcRouter.GET("/plugins/:pluginName", pc.GetUserPluginDetail)
		pcRouter.POST("/plugins/:pluginName/versions", pc.AddVersion)
	}
}

// @Tags	developer
// @Security ApiKeyAuth
// @Summary	add a plugin
// @Produce	json
// @Accept	json
// @Param	data	body 	models.Plugin	true	"plugin data"
// @Success	201	{object}	models.BaseResponse{}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	401 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /developer/plugins [post]
func (pc *PluginsController) Add(c *gin.Context) {
	var p models.Plugin
	if err := c.ShouldBind(&p); err != nil {
		pc.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	userName := pc.GetUserName(c)
	if userName == "" {
		pc.responseBadRequest(c, global.InvalidUserCode, global.InvalidUserMsg)
		return
	}

	err := p.Save(userName)
	if err == nil {
		pc.responseSuccessCreate(c, nil)
	} else {
		pc.responseWithError(c, err)
	}

}

// @Tags	developer
// @Security ApiKeyAuth
// @Summary	Plug-in developers get the list of plug-ins created by themselves
// @Accept  json
// @Produce	json
// @Success	200	{object}	models.BaseResponse{data=[]dbmodels.Plugin}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	403 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /developer/plugins [get]
func (pc *PluginsController) GetUserPlugins(c *gin.Context) {
	userName := pc.GetUserName(c)
	plugins, err := models.GetPluginsByUser(userName)
	if err != nil {
		pc.responseWithError(c, err)
		return
	}
	pc.responseSuccess(c, plugins)
}

// @Tags	developer
// @Summary	Plug-in developers get the details of plug-ins created by themselves
// @Security ApiKeyAuth
// @Accept  json
// @Produce	json
// @Param	pluginName	path	string	true "plugin name"
// @Success	200	{object}	models.BaseResponse{data=dbmodels.Plugin}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	403 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /developer/plugins/{pluginName} [get]
func (pc *PluginsController) GetUserPluginDetail(c *gin.Context) {
	userName := pc.GetUserName(c)
	pluginName := c.Param("pluginName")
	if pluginName == "" {
		pc.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}

	if details, err := models.GetPluginDetails(userName, pluginName); err != nil {
		pc.responseWithError(c, err)
	} else {
		pc.responseSuccess(c, details)
	}

}

// @Tags	developer
// @Summary submit a new version for a plugin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param   pluginName	path	string	true	"the plug-in corresponding to the version "
// @Param	data	body	models.PVersionOPT true	"version info"
// @Success	201	{object}	models.BaseResponse{}	"success result"
// @Failure	400 {object}	models.BaseResponse{}	""
// @Failure	401 {object}	models.BaseResponse{}	""
// @Failure	500	{object}	models.BaseResponse	""
// @Router /developer/plugins/{pluginName}/versions [post]
func (pc *PluginsController) AddVersion(c *gin.Context) {
	pluginName := c.Param("pluginName")
	var pv models.PVersionOPT
	err := c.ShouldBind(&pv)
	if err != nil || pluginName == "" {
		pc.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	uName := pc.GetUserName(c)
	if gErr := pv.AddVersion(pluginName, uName); gErr != nil {
		pc.responseWithError(c, gErr)
	} else {
		pc.responseSuccessCreate(c, nil)
	}

}
