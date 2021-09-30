package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/opensourceways/app-robot-server/config"
	v1 "github.com/opensourceways/app-robot-server/controller/v1"
	_ "github.com/opensourceways/app-robot-server/docs"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/middleware"
)

type IRouter interface {
	RegisterRouter(rg *gin.RouterGroup)
}

func Init() *gin.Engine {
	ginRouter := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpwd", global.StrongPwd)
	}

	if config.Application.RunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginRouter.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": global.NotFoundCode, "msg": global.NotFoundMsg})
	})
	//setting logger and recover middleware
	ginRouter.Use(middleware.Logger(), gin.Recovery())
	//setting swagger doc address
	ginRouter.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	//router group
	v1RG := ginRouter.Group("/v1")
	{
		registerRouters(v1RG, &v1.AccountController{}, &v1.PluginsController{}, &v1.InstanceController{})
	}
	return ginRouter
}

func registerRouters(rg *gin.RouterGroup, routers ...IRouter) {
	for _, v := range routers {
		v.RegisterRouter(rg)
	}
}
