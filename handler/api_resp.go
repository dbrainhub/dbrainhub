package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/gin-gonic/gin"
)

func BadRequestError(format string, a ...interface{}) error {
	return &ApiError{
		Code: http.StatusBadRequest,
		Name: "BadRequestError",
		Msg:  fmt.Sprintf(format, a...),
	}
}

func NewApiError(err error) *ApiError {
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

type ApiError struct {
	Code int32  `json:"code"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func (a *ApiError) Error() string {
	res, _ := json.Marshal(a)
	return string(res)
}

func (a *ApiError) RenderJson(c *gin.Context) {
	statusCode := http.StatusBadRequest

	if a.Code < errors.MinCustomErrorCode {
		statusCode = int(a.Code)
	}

	c.AbortWithStatusJSON(statusCode, a)
}
