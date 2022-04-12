package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type TagItem struct {
	Id        int32  `gorm:"column:id" json:"id"`
	ItemType  string `gorm:"column:item_type" json:"item_type"`
	ItemId    int32  `gorm:"column:item_id" json:"item_id"`
	Tag       string `gorm:"tag" json:"tag"`
	CreatedAt int64  `gorm:"column:ct" json:"created_at"`
	UpdatedAt int64  `gorm:"column:ut" json:"updated_at"`
}

const TableTagItem = "tag_item"

const (
	ItemTypeDBCluster       = "dbcluster"
	ItemTypeDBClusterMember = "dbcluster_member"
)

func GetAllTags(ctx context.Context, db *gorm.DB) ([]string, error) {
	tags := []string{}
	err := db.Table(TableTagItem).Distinct("tag").Find(&tags).Error
	return tags, err
}

func GetTags(ctx context.Context, db *gorm.DB, itemType string, itemId int32) ([]string, error) {
	tags := []string{}
	err := db.Table(TableTagItem).Select("tag").Where("item_type = ? AND item_id = ?", itemType, itemId).Find(&tags).Error
	return tags, err
}

func BatchGetTags(ctx context.Context, db *gorm.DB, itemType string, itemIds []int32) (map[int32][]string, error) {
	type Row struct {
		ItemId int32  `gorm:"column:item_id"`
		Tag    string `gorm:"tag"`
	}
	rows := []Row{}
	err := db.Table(TableTagItem).Select("item_id", "tag").Where("item_type = ? AND item_id in ?", itemType, itemIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	res := map[int32][]string{}
	for _, row := range rows {
		if tags, ok := res[row.ItemId]; ok {
			tags = append(tags, row.Tag)
			res[row.ItemId] = tags
		} else {
			res[row.ItemId] = []string{row.Tag}
		}
	}
	return res, nil
}

func GetItemsByTag(ctx context.Context, db *gorm.DB, tag string, itemType string) ([]int32, error) {
	itemIds := []int32{}
	err := db.Table(TableTagItem).Select("item_id").Where("tag = ? AND item_type = ?", tag, itemType).Find(&itemIds).Error
	return itemIds, err
}

func AlreadyTagged(ctx context.Context, db *gorm.DB, itemType string, itemId int32, tag string) (bool, error) {
	var tagItem TagItem
	err := db.Table(TableTagItem).Where("item_type = ? AND item_id = ? AND tag = ?", itemType, itemId, tag).First(&tagItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func AddTag(ctx context.Context, db *gorm.DB, itemType string, itemId int32, tag string) error {
	now := time.Now().Unix()
	tagItem := TagItem{
		Id:        0,
		ItemType:  itemType,
		ItemId:    itemId,
		Tag:       tag,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return db.Create(&tagItem).Error
}
