package model

import (
	"context"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	DbTypeMySQL = "mysql"
)

type DbCluster struct {
	Id          int32  `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	DbType      string `gorm:"column:db_type" json:"db_type"`
	CreatedAt   int64  `gorm:"column:ct" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:ut" json:"updated_at"`
}

func (cluster *DbCluster) TableName() string {
	return "dbcluster"
}

func GetDbClusters(ctx context.Context, db *gorm.DB, offset int64, limit int64) ([]*DbCluster, error) {
	var clusters []*DbCluster
	err := db.Offset(offset).Limit(limit).Find(&clusters).Error
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func GetDbClusterById(ctx context.Context, db *gorm.DB, id int32) (*DbCluster, error) {
	var cluster DbCluster
	err := db.First(&cluster, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.DbClusterNotFoundById(id)
		}
		return nil, err
	}

	return &cluster, nil
}

func GetDbClusterByName(ctx context.Context, db *gorm.DB, name string) (*DbCluster, error) {
	var cluster DbCluster
	err := db.First(&cluster, "`name` = ?", name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.DbClusterNotFoundByName(name)
		}
		return nil, err
	}
	return &cluster, nil
}

func CreateDbCluster(ctx context.Context, db *gorm.DB, name string, description string, dbType string) (*DbCluster, error) {
	now := time.Now().Unix()
	newCluster := DbCluster{
		// Id:          0,
		Name:        name,
		Description: description,
		DbType:      dbType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := db.Create(&newCluster).Error
	if err != nil {
		return nil, err
	}
	return &newCluster, nil
}
