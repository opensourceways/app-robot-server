package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/models/request"
	"github.com/opensourceways/app-robot-server/service"
)

type AccountController struct {
	baseController
}

func (ac *AccountController) RegisterRouter(rg *gin.RouterGroup) {
	accRouter := rg.Group("/accounts")
	{
		accRouter.POST("/login", ac.Login)
		accRouter.POST("/register", ac.Register)
	}
}

// @Tags Accounts
// @Summary 账号密码登录
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	account		formData	string	true	"账号"
// @Param	password	formData	string	true	"密码"
// @Success 200 {object}	response.BaseResponse{data=response.LoginResult}	"返回token和userId"
// @Failure	400	{object}	response.BaseResponse	"错误返回"
// @Failure	403	{object}	response.BaseResponse	"错误返回"
// @Failure	500	{object}	response.BaseResponse	"错误返回"
// @Router /accounts/login [post]
func (ac *AccountController) Login(c *gin.Context) {
	var param request.Login
	if err := c.ShouldBind(&param); err != nil {
		ac.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	result, err := service.DoLogin(param)
	if err != nil {
		ac.responseWithError(c, err)
		return
	}
	ac.responseSuccess(c, result)
}

// @Tags Accounts
// @Summary 账号注册
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	account		formData	string	true	"邮箱"
// @Param	code		formData	string	true	"验证码"
// @Param 	username 	formData 	string	true	"用户名"
// @Param	password	formData	string	true	"密码"
// @Success 201	{object}	response.BaseResponse		"返回注册状态"
// @Router	/accounts/register	[post]
func (ac *AccountController) Register(c *gin.Context) {
	ac.responseSuccessCreate(c, nil)
}

func (ac *AccountController) FindPassword(c *gin.Context) {

}

func (ac *AccountController) ResetPassword(c *gin.Context) {

}
