package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/app-robot-server/middleware"

	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/logs"
	"github.com/opensourceways/app-robot-server/models"
)

type AccountController struct {
	baseController
}

func (ac *AccountController) RegisterRouter(rg *gin.RouterGroup) {
	accRouter := rg.Group("/accounts")
	{
		accRouter.POST("/login", ac.Login)
		accRouter.POST("/register", ac.Register)
		accRouter.GET("/email", ac.CheckEmail)
		accRouter.GET("/username", ac.CheckUserName)
		accRouter.PUT("/password", middleware.JWTAuth(), ac.ResetPassword)
		accRouter.PATCH("/password", ac.ResetPwdByEmailValidate)
		accRouter.GET("/password", ac.SendResetPasswordEmail)
	}
}

// @Tags Accounts
// @Summary 账号密码登录
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	account		formData	string	true	"账号"
// @Param	password	formData	string	true	"密码"
// @Success 200 {object}	models.BaseResponse{data=models.LoginResult}	"返回token和userId"
// @Failure	400	{object}	models.BaseResponse	"错误返回"
// @Failure	403	{object}	models.BaseResponse	"错误返回"
// @Failure	500	{object}	models.BaseResponse	"错误返回"
// @Router /accounts/login [post]
func (ac *AccountController) Login(c *gin.Context) {
	var param models.Login
	if err := c.ShouldBind(&param); err != nil {
		ac.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	result, err := param.DoLogin()
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
// @Param	email		formData	string	true	"邮箱"
// @Param	code		formData	string	true	"验证码"
// @Param 	username 	formData 	string	true	"用户名(仅字母数字组合)"
// @Param	password	formData	string	true	"密码（大小字母数字特殊字符满足三种）"
// @Success 201	{object}	models.BaseResponse		"返回成功状态"
// @Failure	400 {object}	models.BaseResponse		"返回失败状态"
// @Failure	500 {object}	models.BaseResponse		"返回失败状态"
// @Router	/accounts/register	[post]
func (ac *AccountController) Register(c *gin.Context) {
	var param models.Register
	if err := c.ShouldBind(&param); err != nil {
		logs.Logger.Error(err)
		ac.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	if rErr := param.DoRegister(); rErr != nil {
		ac.responseWithError(c, rErr)
		return
	}
	ac.responseSuccessCreate(c, nil)
}

// @Tags Accounts
// @Summary 发送找回密码邮件
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	email		formData	string	true	"邮箱"
// @Success 201	{object}	models.BaseResponse		"邮件发送成功"
// @Failure	400 {object}	models.BaseResponse		"邮箱校验失败"
// @Failure	500 {object}	models.BaseResponse		"返回失败状态"
// @Router /accounts/password [get]
func (ac *AccountController) SendResetPasswordEmail(c *gin.Context) {
	//TODO:implement
	ac.responseServerError(c)
}

// @Tags Accounts
// @Summary 邮件验证重置密码
// @Produce json
// @Accept x-www-form-urlencoded
// @param	key			formData	string	true	"the key"
// @param	email		formData	string	true	"email"
// @param	password 	formData	string	true	"new password"
// @Success 201	{object}	models.BaseResponse		"邮件发送成功"
// @Failure	403 {object}	models.BaseResponse		"验证失败"
// @Failure	400 {object}	models.BaseResponse		"重置失败"
// @Failure	500 {object}	models.BaseResponse		"返回失败状态"
// @Router	/accounts/password [patch]
func (ac *AccountController) ResetPwdByEmailValidate(c *gin.Context) {
	//TODO:implement
	ac.responseServerError(c)
}

// @Tags Accounts
// @Summary 重置密码
// @Security ApiKeyAuth
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	userID		formData	string	true	"用户ID"
// @param	password	formData	string	true	"旧密码"
// @param	newPassword	formData	string	true	"新密码"
// @Success	201	{object}	models.BaseResponse		"修改成功"
// @Failure 401 {object}	models.BaseResponse		"token验证失败"
// @Failure 400 {object}	models.BaseResponse		"修改失败"
// @Failure 500	{object}	models.BaseResponse		"服务器错误"
// @Router	/accounts/password [put]
func (ac *AccountController) ResetPassword(c *gin.Context) {
	var p models.ResetPwdBM
	if err := c.ShouldBind(&p); err != nil {
		ac.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}

	if err := p.DoResetPwd(); err == nil {
		ac.responseSuccess(c, nil)
	} else {
		ac.responseWithError(c, err)
	}

}

// @Tags Accounts
// @Summary 账号注册获取验证码
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	email		query	string	true	"邮箱"
// @Success 200	{object}	models.BaseResponse{}		"邮件发送成功"
// @Router	/accounts/code	[get]
func (ac *AccountController) RegisterCode(c *gin.Context) {
	//TODO:implement
	ac.responseServerError(c)
}

// @Tags Accounts
// @Summary 邮箱是否已注册
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	email		query	string	true	"邮箱"
// @Success 200	{object}	models.BaseResponse		"邮箱可用"
// @Failure	400 {object}	models.BaseResponse		"邮箱已注册"
// @Failure	500 {object}	models.BaseResponse		"返回失败状态"
// @Router	/accounts/email	[get]
func (ac *AccountController) CheckEmail(c *gin.Context) {
	var p models.EmailBinding
	if err := c.ShouldBindQuery(&p); err != nil {
		logs.Logger.Error(err)
		ac.responseBadRequest(c, global.IllegalInputErrCode, global.IllegalInputErrMsg)
		return
	}
	registered, err := p.EmailRegistered()
	if err != nil {
		logs.Logger.Error(err)
		ac.responseWithError(c, global.NewResponseSystemError())
		return
	}

	if !registered {
		ac.responseSuccess(c, nil)
	} else {
		ac.responseBadRequest(c, global.RegisteredErrCode, global.RegisteredErrMsg)
	}
}

// @Tags Accounts
// @Summary 用户名是否可用
// @Produce json
// @Accept x-www-form-urlencoded
// @Param	username		query	string	true	"用户名"
// @Success 200	{object}	models.BaseResponse		"用户名可用"
// @Failure	500 {object}	models.BaseResponse		"用户名已注册"
// @Router	/accounts/username	[get]
func (ac *AccountController) CheckUserName(c *gin.Context) {
	userName := c.Query("username")
	if userName == "" {
		ac.responseWithError(c, global.NewResponseSystemError())
		return
	}
	exist, err := dbmodels.GetDB().LoginNameExist(userName)
	if err != nil {
		logs.Logger.Error(err)
		ac.responseWithError(c, global.NewResponseSystemError())
		return
	}
	if exist {
		ac.responseBadRequest(c, global.RegisteredErrCode, global.RegisteredErrMsg)
	} else {
		ac.responseSuccess(c, nil)
	}

}
