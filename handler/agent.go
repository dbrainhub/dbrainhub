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

func autoAdjustLimitAndOffset(limit, offset int32) (int32, int32) {
	if offset <= 0 {
		offset = DefaultOffset
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return limit, offset
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
