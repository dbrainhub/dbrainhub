package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

const (
	DefaultOffset = 0
	DefaultLimit  = 20
	MaxLimit      = 100
)

type OffsetLimit struct {
	Offset int64 `json:"offset" form:"offset" uri:"offset"`
	Limit  int64 `json:"limit" form:"offset" uri:"limit"`
}

func (ol *OffsetLimit) AutoAdjust() {
	if ol.Offset <= 0 {
		ol.Offset = DefaultOffset
	}
	if ol.Limit <= 0 {
		ol.Limit = DefaultLimit
	}
	if ol.Limit > MaxLimit {
		ol.Limit = MaxLimit
	}
}

func SayHello(c *gin.Context) (interface{}, error) {
	var req *api.HelloWorldRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.NewController().SayHello(c, req)

}
