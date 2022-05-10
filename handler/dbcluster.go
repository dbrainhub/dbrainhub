package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func GetDbClusters(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var req api.GetDBClustersRequest
	if err := bindQuery(c, &req); err != nil {
		return nil, err
	}
	req.Limit, req.Offset = autoAdjustLimitAndOffset(req.Limit, req.Offset)

	return controller.GetDbClusters(c, currUser, int(req.Offset), int(req.Limit))
}

func CreateDbCluster(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var req api.NewDBClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, BadRequestError("json unsharmal err: %v", err)
	}
	return controller.CreateDbCluster(c, currUser, &req)
}
