package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context) (interface{}, error) {
	var req *api.HelloWorldRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.SayHello(c, req)
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
