package controller

import (
	"context"
	"fmt"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
)

func GetDbClusters(ctx context.Context, currUser *model.User, offset, limit int) (*api.GetDBClustersResponse, error) {
	db := model.GetDB(ctx)
	clusters, err := model.GetDbClusters(ctx, db, offset, limit)
	if err != nil {
		return nil, err
	}
	return &api.GetDBClustersResponse{
		Dbclusters: toDBClusters(clusters),
	}, nil
}

func CreateDbCluster(ctx context.Context, currUser *model.User, params *api.NewDBClusterRequest) (*api.DBCluster, error) {
	if err := validateNewDBClusterRequest(params); err != nil {
		return nil, err
	}
	db := model.GetDB(ctx)
	cluster, err := model.CreateDbCluster(ctx, db, params.Name, params.Description, params.DbType)
	if err != nil {
		return nil, err
	}

	err = model.BatchAssignMembersToCluster(ctx, db, params.MemberIds, cluster.Id)
	if err != nil {
		return nil, err
	}
	return toDBCluster(cluster), nil
}

func validateNewDBClusterRequest(params *api.NewDBClusterRequest) error {
	if params.DbType != model.DbTypeMySQL {
		return errors.InvalidDbType(fmt.Sprintf("invalid db_type: %s", params.DbType))
	}
	return nil
}

func toDBClusters(clusters []*model.DbCluster) []*api.DBCluster {
	var res []*api.DBCluster
	for _, cluster := range clusters {
		res = append(res, toDBCluster(cluster))
	}
	return res
}

func toDBCluster(cluster *model.DbCluster) *api.DBCluster {
	return &api.DBCluster{
		Id:          cluster.Id,
		Name:        cluster.Name,
		Description: cluster.Description,
		Dbtype:      cluster.DbType,
		CreatedAt:   cluster.CreatedAt,
		UpdatedAt:   cluster.UpdatedAt,
	}
}
