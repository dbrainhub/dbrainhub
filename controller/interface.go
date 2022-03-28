package controller

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	SayHello(*gin.Context, *api.HelloWorldRequest) (*api.HelloWorldResponse, error)
}

func NewController() Controller {
	return &DefaultController{}
}
