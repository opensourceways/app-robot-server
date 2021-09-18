package models

import (
	"github.com/google/uuid"

	"github.com/opensourceways/app-robot-server/dbmodels"
	"github.com/opensourceways/app-robot-server/global"
	"github.com/opensourceways/app-robot-server/logs"
	"github.com/opensourceways/app-robot-server/utils"
)

type LoginResult struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type Login struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (l *Login) DoLogin() (*LoginResult, global.Error) {
	user, dErr := dbmodels.GetDB().GetUser(l.Account, utils.MD5Encoding(l.Password))
	var repErr global.ResponseError
	if dErr != nil {
		if dErr.IsErrorOf(dbmodels.ErrNoDBRecord) {
			repErr = global.ResponseError{ErrCode: global.AccountPasswordErrCode, Reason: global.AccountPasswordErrMsg}
		} else {
			repErr = global.NewResponseSystemError()
		}
		return nil, repErr
	}

	token, err := genToken(user.ID)
	if err != nil {
		logs.Logger.Error(err)
		return nil, global.ResponseError{ErrCode: global.SystemErrorCode, Reason: global.ServerErrorMsg}
	}
	//TODO: serverTheTokenToCache(token)
	return &LoginResult{
		UserID: user.ID,
		Token:  token,
	}, nil
}

type EmailBinding struct {
	Email string `form:"email" query:"email" binding:"required,email"`
}

func (eb EmailBinding) EmailRegistered() (bool, error) {
	return dbmodels.GetDB().EmailExist(eb.Email)
}

type Register struct {
	EmailBinding
	UserName string `form:"username" binding:"required,alphanum"`
	Password string `form:"password" binding:"required,strongpwd"`
	Code     string `form:"code" binding:"required"`
}

func (r Register) DoRegister() global.Error {
	userId := uuid.NewString()
	dUser := dbmodels.User{
		ID:        userId,
		LoginName: r.UserName,
		Email:     r.Email,
		Password:  utils.MD5Encoding(r.Password),
	}
	if dErr := dbmodels.GetDB().AddUser(dUser); dErr != nil {
		if dErr.IsErrorOf(dbmodels.ErrNoDBRecord) {
			return global.ResponseError{ErrCode: global.RegisteredErrCode, Reason: global.RegisteredErrMsg}
		} else {
			return global.NewResponseSystemError()
		}
	}
	return nil
}

type ResetPwdBM struct {
	UserID      string `form:"userID" binding:"required"`
	Password    string `form:"password" binding:"required"`
	NewPassword string `form:"newPassword" binding:"required,strongpwd"`
}

//DoResetPwd reset user password by user id and password
func (rpb ResetPwdBM) DoResetPwd() global.Error {
	oldPwd := utils.MD5Encoding(rpb.Password)
	newPwd := utils.MD5Encoding(rpb.NewPassword)
	err := dbmodels.GetDB().UpdatePassword(rpb.UserID, oldPwd, newPwd)
	if err == nil {
		return nil
	}
	if err.IsErrorOf(dbmodels.ErrNoDBRecord) {
		return global.ResponseError{ErrCode: global.NoUserOrPasswordCode, Reason: global.NoUserOrPasswordMsg}
	}
	return global.NewResponseSystemError()
}
