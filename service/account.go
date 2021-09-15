package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/opensourceway/app-robot-server/config"
	"github.com/opensourceway/app-robot-server/global"
	"github.com/opensourceway/app-robot-server/logs"
	"github.com/opensourceway/app-robot-server/models"
	"github.com/opensourceway/app-robot-server/models/request"
	"github.com/opensourceway/app-robot-server/models/response"
)

func DoLogin(params request.Login) (*response.LoginResult, global.Error) {
	//TODO : check login illegal
	userId := "zggdhh677"
	token, err := genToken("zggdhh677")
	if err != nil {
		logs.Logger.Error(err)
		return nil, global.ResponseError{ErrCode: global.SystemErrorCode, Reason: global.ServerErrorMsg}
	}
	//TODO: serverTheTokenToCache(token)
	return &response.LoginResult{
		UserID: userId,
		Token:  token,
	}, nil

}

func genToken(userId string) (string, error) {
	j := models.NewJwt()
	claims := models.CustomClaims{
		ID: userId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
			ExpiresAt: config.Application.Jwt.TokenExpiration * 60 * 60,
			Issuer:    "app-rabot-server",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
