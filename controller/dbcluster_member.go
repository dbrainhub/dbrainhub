package controller

import (
	"context"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/model"
)

func GetToAssignDbClusterMembers(ctx context.Context, currUser *model.User, dbtype, env string, ipPrefix *string, limit, offset int32) (*api.GetToAssignDbClusterMembersResponse, error) {
	db := model.GetDB(ctx)
	members, err := model.GetToAssignClusterMembers(ctx, db, dbtype, env, ipPrefix, int(limit), int(offset))
	if err != nil {
		return nil, err
	}
	clusters, err := model.GetDbClusterByIds(ctx, db, members.UniqValidClusterIds())
	if err != nil {
		return nil, err
	}
	return &api.GetToAssignDbClusterMembersResponse{
		Members: assignDBClusterMembersExtra(toDBClusterMembers(members), clusters),
	}, nil
}

func GetClusterMembers(ctx context.Context, currUser *model.User, clusterId int32) (*api.GetDbClusterMembersResponse, error) {
	db := model.GetDB(ctx)
	cluster, err := model.GetDbClusterById(ctx, db, clusterId)
	if err != nil {
		return nil, err
	}

	if err := UserCanAccessDbCluster(currUser, cluster); err != nil {
		return nil, err
	}
	members, err := model.GetDbClusterMembers(ctx, db, clusterId)
	if err != nil {
		return nil, err
	}
	return &api.GetDbClusterMembersResponse{
		Members: toDBClusterMembers(members),
	}, nil
}

func AssignClusterMembers(ctx context.Context, currUser *model.User, clusterId int32, memberIds []int32) error {
	db := model.GetDB(ctx)
	cluster, err := model.GetDbClusterById(ctx, db, clusterId)
	if err != nil {
		return err
	}

	if err := UserCanAccessDbCluster(currUser, cluster); err != nil {
		return err
	}

	// 新的 Assign 会覆盖掉之前的。允许向一个集群中添加多次实例。
	err = model.BatchAssignMembersToCluster(ctx, db, memberIds, clusterId)
	if err != nil {
		return err
	}

	return nil
}

// UserCanAccessDbCluster checks if user has the permission to access db cluster.
// placeholder for now.
func UserCanAccessDbCluster(user *model.User, cluster *model.DbCluster) error {
	return nil
}

func assignDBClusterMembersExtra(members []*api.DBClusterMember, clusters []*model.DbCluster) []*api.DBClusterMember {
	for _, member := range members {
		assignDBClusterMemberExtra(member, clusters)
	}
	return members
}

func assignDBClusterMemberExtra(member *api.DBClusterMember, clusters []*model.DbCluster) {
	if member.ClusterId == 0 {
		return
	}
	for _, cluster := range clusters {
		if cluster.Id == member.ClusterId {
			member.Extra.ClusterName = cluster.Name
		}
	}
}

func toDBClusterMembers(members []*model.DbClusterMember) []*api.DBClusterMember {
	var res []*api.DBClusterMember
	for _, member := range members {
		res = append(res, toDBClusterMember(member))
	}
	return res
}

func toDBClusterMember(member *model.DbClusterMember) *api.DBClusterMember {
	return &api.DBClusterMember{
		Id:        member.Id,
		ClusterId: member.ClusterId,
		Hostname:  member.Hostname,
		DbType:    member.DbType,
		DbVersion: member.DbVersion,
		Role:      member.Role,
		Ipaddr:    member.IPaddr,
		Port:      int32(member.Port),
		Os:        member.Os,
		OsVersion: member.OsVersion,
		HostType:  member.HostType,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}
