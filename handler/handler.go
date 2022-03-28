package handler

import (
	"fmt"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context) {
	var req *api.HelloWorldRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		BadRequestResp(c, fmt.Sprintf("input unsharmal err: %v", err))
		return
	}
	res, err := controller.NewController().SayHello(c, req)
	if err != nil {
		FailResp(c, err)
		return
	}
	SuccessResp(c, res)
}
