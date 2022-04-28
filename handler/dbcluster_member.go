package handler

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func GetToAssignedDbClusterMembers(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var req api.GetToAssignDbClusterMembersRequest
	if err := c.BindQuery(&req); err != nil {
		return nil, err
	}
	req.Limit, req.Offset = autoAdjustLimitAndOffset(req.Limit, req.Offset)

	return controller.GetToAssignDbClusterMembers(c, currUser, req.DbType, req.Env, req.IpPrefix, req.Limit, req.Offset)
}

func GetDbClusterMembers(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var params struct {
		ClusterId int32 `uri:"clusterId" binding:"required"`
	}
	if err := c.BindUri(&params); err != nil {
		return nil, err
	}

	return controller.GetClusterMembers(c, currUser, params.ClusterId)
}

func AssignDbClusterMembers(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var pClusterId struct {
		ClusterId int32 `uri:"clusterId" binding:"required"`
	}
	if err := c.BindUri(&pClusterId); err != nil {
		return nil, err
	}

	var req api.AssignDbClusterMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, BadRequestError("json unsharmal err: %v", err)
	}

	err = controller.AssignClusterMembers(c, currUser, pClusterId.ClusterId, req.MemberIds)
	if err != nil {
		return nil, err
	}
	return &api.AssignDbClusterMembersResponse{}, nil
}
