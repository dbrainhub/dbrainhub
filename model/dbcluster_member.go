package model

import "context"

type DbClusterMember struct {
	Id        int32  `gorm:"column:id"`
	ClusterId int32  `gorm:"column:cluster_id"`
	IP        string `gorm:"column:ip"`
	CreatedAt int64  `gorm:"column:ct"`
	UpdatedAt int64  `gorm:"column:ut"`
}

func GetDbClusterMembers(ctx context.Context, clusterId int32, offset int64, limit int64) ([]*DbClusterMember, error) {
	var members []*DbClusterMember
	err := db.Where("`cluster_id` = ?").Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
