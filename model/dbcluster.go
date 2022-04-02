package model

import (
	"context"
	"github.com/jinzhu/gorm"
)

type DbCluster struct {
	Id        int32  `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	Type      string `gorm:"column:type"`
	CreatedAt int64  `gorm:"column:ct"`
	UpdatedAt int64  `gorm:"column:ut"`
}

func GetDbClusters(ctx context.Context, db *gorm.DB, offset int64, limit int64) ([]*DbCluster, error) {
	var clusters []*DbCluster
	err := db.Offset(offset).Limit(limit).Find(&clusters).Error
	if err != nil {
		return nil, err
	}
	return clusters, nil
}
