package controller

import (
	"context"

	"github.com/dbrainhub/dbrainhub/api"
)

type DefaultController struct{}

func (d *DefaultController) Hello(context.Context, *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	return &api.HelloWorldResponse{
		Pang: "hello,world!",
	}, nil
}
