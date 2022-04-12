package controller

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context, req *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	return &api.HelloWorldResponse{
		Pang: "hello,world!",
	}, nil
}

func Heartbeat(c *gin.Context, req *api.HeartbeatRequest) (*api.HeartbeatResponse, error) {
	logger.Infof("receive heartbeat from IP:%s, port: %d, req: %#v\n", req.AgentInfo.Localip, req.DbInfo.Port, req)
	return &api.HeartbeatResponse{}, nil
}

func Report(c *gin.Context, req *api.StartupReportRequest) (*api.StartupReportResponse, error) {
	logger.Infof("receive reporter from IP:%s, port: %d, req: %#v\n", req.IpAddr, req.Port, req)
	return &api.StartupReportResponse{}, nil
}
