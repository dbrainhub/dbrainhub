package handler

import (
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func GetUnassignedDbClusterMembers(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var p OffsetLimit
	if err := c.BindQuery(&p); err != nil {
		return nil, err
	}
	p.AutoAdjust()

	members, err := controller.GetUnassignedDbClusterMembers(c, currUser, p.Offset, p.Limit)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"members": members}, nil
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

	members, err := controller.GetClusterMembers(c, currUser, params.ClusterId)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"members": members}, nil
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

	var pMemberIds struct {
		MemberIds []int32 `json:"member_ids"`
	}
	if err := c.ShouldBindJSON(&pMemberIds); err != nil {
		return nil, BadRequestError("json unsharmal err: %v", err)
	}

	err = controller.AssignClusterMembers(c, currUser, pClusterId.ClusterId, pMemberIds.MemberIds)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
