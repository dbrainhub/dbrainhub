package handler

import (
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func GetDbClusters(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var p OffsetLimit
	if err := c.BindQuery(&p); err != nil {
		return nil, err
	}
	p.AutoAdjust()

	clusters, err := controller.GetDbClusters(c, currUser, p.Offset, p.Limit)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"dbclusters": clusters}, nil
}

func CreateDbCluster(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	params := controller.CreateDbClusterParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		return nil, BadRequestError("json unsharmal err: %v", err)
	}
	return controller.CreateDbCluster(c, currUser, &params)
}
