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
	Offset int `json:"offset" form:"offset" uri:"offset"`
	Limit  int `json:"limit" form:"offset" uri:"limit"`
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

func Heartbeat(c *gin.Context) (interface{}, error) {
	var req *api.HeartbeatRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.Heartbeat(c, req)
}

func Report(c *gin.Context) (interface{}, error) {
	var req *api.StartupReportRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.Report(c, req)
}
