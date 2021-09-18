package v1

import (
	"github.com/opensourceways/app-robot-server/models"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/app-robot-server/global"
)

type baseController struct {
}

func (bc *baseController) responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.BaseResponse{Code: global.SuccessCode, Msg: global.SuccessMsg, Data: data})
}

func (bc *baseController) responseSuccessCreate(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.BaseResponse{Code: global.SuccessCode, Msg: global.SuccessMsg, Data: data})
}

func (bc *baseController) responseServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, models.BaseResponse{Code: global.SystemErrorCode, Msg: global.ServerErrorMsg})
}

func (bc *baseController) responseBadRequest(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusBadRequest, models.BaseResponse{Code: code, Msg: msg})
}

func (bc *baseController) responseForbidden(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusForbidden, models.BaseResponse{Code: code, Msg: msg})
}

func (bc *baseController) responseServerErrorWithCode(c *gin.Context, statusCode int, err global.Error) {
	c.JSON(statusCode, models.BaseResponse{Code: err.Code(), Msg: err.Msg()})
}

func (bc *baseController) responseWithError(c *gin.Context, err global.Error) {
	if err == nil {
		bc.responseServerError(c)
		return
	}
	statusCode := http.StatusBadRequest
	switch err.Code() {
	case global.SystemErrorCode, global.UnknownCacheErrorCode, global.UnknownDBErrorCode:
		statusCode = http.StatusInternalServerError
	case global.UnauthorizedCode:
		statusCode = http.StatusUnauthorized
	case global.NotFoundCode:
		statusCode = http.StatusNotFound
	}
	bc.responseServerErrorWithCode(c, statusCode, err)
}
