package controller

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	SayHello(*gin.Context, *api.HelloWorldRequest) (*api.HelloWorldResponse, error)
	DbRainhubOutput(*gin.Context, DbRainhubRequest) (*DbRainhubResponse, error)
}

func NewController() Controller {
	return &DefaultController{}
}
