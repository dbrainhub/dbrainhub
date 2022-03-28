package controller

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/gin-gonic/gin"
)

type DefaultController struct{}

func (d *DefaultController) SayHello(c *gin.Context, req *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	return &api.HelloWorldResponse{
		Pang: "hello,world!",
	}, nil
}
