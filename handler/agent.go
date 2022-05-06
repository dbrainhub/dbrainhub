package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	DefaultOffset = 0
	DefaultLimit  = 20
	MaxLimit      = 100
)

func autoAdjustLimitAndOffset(limit, offset int32) (int32, int32) {
	if offset < 0 {
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

// 将 GET 请求参数映射到结构体，需要结构体字段包含 form tag。proto 不支持添加自定义 tag。
func bindQuery(c *gin.Context, obj interface{}) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}
	return binding.MapFormWithTag(obj, c.Request.Form, "json") // 使用 json tag
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
