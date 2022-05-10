package handler

import (
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func DbRainhubOutput(c *gin.Context) (interface{}, error) {
	var req controller.DbRainhubRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.DbRainhubOutput(c, req)
}
