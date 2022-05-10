package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func DbRainhubSearchMemberLogsWithCount(c *gin.Context) (interface{}, error) {
	var req api.SearchMemberLogCountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, BadRequestError("input unsharmal err: %v", err)
	}
	return controller.DbRainhubSearchMemberLogsWithCount(c, &req)
}
