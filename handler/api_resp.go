package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/gin-gonic/gin"
)

const (
	RetSuccess = 1
	RetFailure = -1
)

func BadRequestResp(c *gin.Context, msg string) {
	resp := &ApiError{
		Code: http.StatusBadRequest,
		Name: "BadRequestError",
		Msg:  msg,
	}
	resp.RenderJson(c)
}

func FailResp(c *gin.Context, err error) {
	resp := newApiError(err)
	resp.RenderJson(c)
}

func SuccessResp(c *gin.Context, data interface{}) {
	resp := &ApiSucc{
		Ret:  RetSuccess,
		Data: data,
	}
	resp.RenderJson(c)
}

type (
	ApiSucc struct {
		Ret  int32       `json:"ret"`
		Data interface{} `json:"data"`
	}

	ApiError struct {
		Ret  int32  `json:"ret"`
		Code int32  `json:"code"`
		Name string `json:"name"`
		Msg  string `json:"msg"`
	}
)

func (a *ApiSucc) RenderJson(g *gin.Context) {
	statusCode := http.StatusOK
	g.JSON(statusCode, a)
}

func (a *ApiError) Error() string {
	res, _ := json.Marshal(a)
	return string(res)
}

func (a *ApiError) RenderJson(c *gin.Context) {
	statusCode := http.StatusOK
	if a.Code == http.StatusInternalServerError || a.Code == http.StatusBadRequest {
		statusCode = int(a.Code)
	}

	a.Ret = RetFailure
	c.AbortWithStatusJSON(statusCode, a)
}

func newApiError(err error) *ApiError {
	switch err.(type) {
	case *ApiError:
		return err.(*ApiError)
	case *errors.ErrInfo:
		mse := err.(*errors.ErrInfo)
		return &ApiError{
			Code: mse.Code,
			Name: mse.Name,
			Msg:  mse.Msg,
		}
	default:
		return &ApiError{
			Code: http.StatusInternalServerError,
			Name: "InternalServerError",
			Msg:  err.Error(),
		}
	}
}
