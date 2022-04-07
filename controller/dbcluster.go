package controller

import (
	"context"
	"fmt"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
)

func GetDbClusters(ctx context.Context, currUser *model.User, offset, limit int64) ([]*model.DbCluster, error) {
	db := model.GetDB(ctx)
	return model.GetDbClusters(ctx, db, offset, limit)
}

type CreateDbClusterParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DbType      string `json:"db_type"`
}

func (cdp *CreateDbClusterParams) Validate() error {
	if cdp.DbType != model.DbTypeMySQL {
		return errors.InvalidDbType(fmt.Sprintf("invalid db_type: %s", cdp.DbType))
	}
	return nil
}

func CreateDbCluster(ctx context.Context, currUser *model.User, params *CreateDbClusterParams) (*model.DbCluster, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	db := model.GetDB(ctx)
	return model.CreateDbCluster(ctx, db, params.Name, params.Description, params.DbType)
}
