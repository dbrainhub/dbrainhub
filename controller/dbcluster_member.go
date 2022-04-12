package controller

import (
	"context"
	"github.com/dbrainhub/dbrainhub/model"
)

func GetUnassignedDbClusterMembers(ctx context.Context, currUser *model.User, offset, limit int) ([]*model.DbClusterMember, error) {
	db := model.GetDB(ctx)
	return model.GetUnassignedClusterMembers(ctx, db, offset, limit)
}

func GetClusterMembers(ctx context.Context, currUser *model.User, clusterId int32) ([]*model.DbClusterMember, error) {
	db := model.GetDB(ctx)
	cluster, err := model.GetDbClusterById(ctx, db, clusterId)
	if err != nil {
		return nil, err
	}

	if err := UserCanAccessDbCluster(currUser, cluster); err != nil {
		return nil, err
	}
	return model.GetDbClusterMembers(ctx, db, clusterId)
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

	// 先查到 clusterId 的所有 member
	// 并根据传入的 memberIds，将原来的 member 分为 3 类分别处理: 1)新增到 cluster 中; 2)从 cluster 中移除; 3)不变
	OldMembers, err := model.GetDbClusterMembers(ctx, db, clusterId)
	if err != nil {
		return err
	}
	oldMemberIds := map[int32]struct{}{}
	for _, m := range OldMembers {
		oldMemberIds[m.Id] = struct{}{}
	}

	memberIdMap := map[int32]struct{}{}
	for _, mid := range memberIds {
		memberIdMap[mid] = struct{}{}
	}

	toAdd := []int32{}
	toDel := []int32{}
	for mid := range memberIdMap {
		if _, ok := oldMemberIds[mid]; !ok {
			toAdd = append(toAdd, mid)
		}
	}
	for mid := range oldMemberIds {
		if _, ok := memberIdMap[mid]; !ok {
			toDel = append(toDel, mid)
		}
	}

	err = model.BatchAssignMembersToCluster(ctx, db, toAdd, clusterId)
	if err != nil {
		return err
	}
	err = model.BatchUnassignClusterMembers(ctx, db, toDel)
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
