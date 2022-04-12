package controller

import (
	"context"
	"fmt"

	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
)

func AddTag(ctx context.Context, currUser *model.User, itemType string, itemId int32, tag string) error {
	db := model.GetDB(ctx)
	switch itemType {
	case model.ItemTypeDBCluster:
		if _, err := model.GetDbClusterById(ctx, db, itemId); err != nil {
			return err
		}
	case model.ItemTypeDBClusterMember:
		if _, err := model.GetDbClusterMemberById(ctx, db, itemId); err != nil {
			return err
		}
	default:
		msg := fmt.Sprintf("invalid item_type: %s", itemType)
		return errors.InvalidItemType(msg)
	}

	alreadTagged, err := model.AlreadyTagged(ctx, db, itemType, itemId, tag)
	if err != nil {
		return err
	}
	if !alreadTagged {
		return model.AddTag(ctx, db, itemType, itemId, tag)
	} else {
		return nil
	}
}

func GetAllTags(ctx context.Context, currUser *model.User) ([]string, error) {
	return model.GetAllTags(ctx, model.GetDB(ctx))
}
