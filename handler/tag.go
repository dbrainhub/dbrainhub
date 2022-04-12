package handler

import (
	"github.com/dbrainhub/dbrainhub/controller"
	"github.com/gin-gonic/gin"
)

func AddTag(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	var params struct {
		ItemType string `json:"item_type"`
		ItemId   int32  `json:"item_id"`
		Tag      string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		return nil, BadRequestError("json unsharmal err: %v", err)
	}

	err = controller.AddTag(c, currUser, params.ItemType, params.ItemId, params.Tag)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func GetAllTags(c *gin.Context) (interface{}, error) {
	currUser, err := Authenticate(c)
	if err != nil {
		return nil, err
	}

	tags, err := controller.GetAllTags(c, currUser)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"tags": tags,
	}
	return res, nil
}
