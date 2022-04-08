package model

import (
	"context"
	"github.com/jinzhu/gorm"
	"time"
)

type DbClusterMember struct {
	Id        int32  `gorm:"column:id" json:"id"`
	ClusterId int32  `gorm:"column:cluster_id" json:"cluster_id"`
	Hostname  string `gorm:"column:hostname" json:"hostname"`
	DbType    string `gorm:"column:db_type" json:"db_type"`
	DbVersion string `gorm:"column:db_version" json:"db_version"`
	Role      int32  `gorm:"column:role" json:"role"`
	IPaddr    string `gorm:"column:ipaddr" json:"ipaddr"`
	Port      int16  `gorm:"column:port" json:"port"`
	Os        string `gorm:"column:os" json:"os"`
	OsVersion string `gorm:"column:os_version" json:"os_version"`
	HostType  int32  `gorm:"column:host_type" json:"host_type"`
	CreatedAt int64  `gorm:"column:ct" json:"created_at"`
	UpdatedAt int64  `gorm:"column:ut" json:"updated_at"`
}

func (cm *DbClusterMember) TableName() string {
	return "dbcluster_member"
}

func GetDbClusterMembers(ctx context.Context, db *gorm.DB, clusterId int32) ([]*DbClusterMember, error) {
	var members []*DbClusterMember
	err := db.Where("`cluster_id` = ?", clusterId).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func GetUnassignedClusterMembers(ctx context.Context, db *gorm.DB, offset int64, limit int64) ([]*DbClusterMember, error) {
	var members []*DbClusterMember
	err := db.Where("`cluster_id` = ?", 0).Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func GetDbClusterMemberByIpAndPort(ctx context.Context, db *gorm.DB, ipAddr string, port int16) (*DbClusterMember, error) {
	var member DbClusterMember
	err := db.First(&member, "`ipaddr` = ? AND `port` = ?", ipAddr, port).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

type NewDbClusterMember struct {
	// Id        int32  `gorm:"column:id"`
	ClusterId int32  `gorm:"column:cluster_id" json:"cluster_id"`
	Hostname  string `gorm:"column:hostname" json:"hostname"`
	DbType    string `gorm:"column:db_type" json:"db_type"`
	DbVersion string `gorm:"column:db_version" json:"db_version"`
	Role      int32  `gorm:"column:role" json:"role"`
	IPaddr    string `gorm:"column:ipaddr" json:"ipaddr"`
	Port      int16  `gorm:"column:port" json:"port"`
	Os        string `gorm:"column:os" json:"os"`
	OsVersion string `gorm:"column:os_version" json:"os_version"`
	HostType  int32  `gorm:"column:host_type" json:"host_type"`
	// CreatedAt int64  `gorm:"column:ct"`
	// UpdatedAt int64  `gorm:"column:ut"`
}

func CreateDbClusterMember(ctx context.Context, db *gorm.DB, params *NewDbClusterMember) (*DbClusterMember, error) {
	now := time.Now().Unix()
	newMember := DbClusterMember{
		// Id:        0,
		ClusterId: params.ClusterId,
		Hostname:  params.Hostname,
		DbType:    params.DbType,
		DbVersion: params.DbVersion,
		Role:      params.Role,
		IPaddr:    params.IPaddr,
		Port:      params.Port,
		Os:        params.Os,
		OsVersion: params.OsVersion,
		HostType:  params.HostType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := db.Create(&newMember).Error
	if err != nil {
		return nil, err
	}
	return &newMember, nil
}

type UpdateDbClusterMemberParams struct {
	Id        int32   `gorm:"column:id" json:"id"`
	ClusterId *int32  `gorm:"column:cluster_id" json:"cluster_id"`
	Hostname  *string `gorm:"column:hostname" json:"hostname"`
	DbType    *string `gorm:"column:db_type" json:"db_type"`
	DbVersion *string `gorm:"column:db_version" json:"db_version"`
	Role      *int32  `gorm:"column:role" json:"role"`
	IPaddr    *string `gorm:"column:ipaddr" json:"ipaddr"`
	Port      *int16  `gorm:"column:port" json:"port"`
	Os        *string `gorm:"column:os" json:"os"`
	OsVersion *string `gorm:"column:os_version" json:"os_version"`
	HostType  *int32  `gorm:"column:host_type" json:"host_type"`
}

func UpdateDbClusterMember(ctx context.Context, db *gorm.DB, params *UpdateDbClusterMemberParams) error {
	mp := map[string]interface{}{}
	if params.ClusterId != nil {
		mp["cluster_id"] = *params.ClusterId
	}
	if params.Hostname != nil {
		mp["hostname"] = *params.Hostname
	}
	if params.DbType != nil {
		mp["db_type"] = *params.DbType
	}
	if params.DbVersion != nil {
		mp["db_version"] = *params.DbVersion
	}
	if params.Role != nil {
		mp["role"] = *params.Role
	}
	if params.IPaddr != nil {
		mp["ipaddr"] = *params.IPaddr
	}
	if params.Port != nil {
		mp["port"] = *params.Port
	}
	if params.Os != nil {
		mp["os"] = *params.Os
	}
	if params.OsVersion != nil {
		mp["os_version"] = *params.OsVersion
	}
	if params.HostType != nil {
		mp["host_type"] = *params.HostType
	}

	if len(mp) == 0 { // nothing to do
		return nil
	}

	mp["ut"] = time.Now().Unix()
	return db.Model(DbClusterMember{}).Where("`id` = ?", params.Id).Update(mp).Error
}

func DeleteDbClusterMemberById(ctx context.Context, db *gorm.DB, id int32) error {
	return db.Delete(DbClusterMember{}, "`id` = ?", id).Error
}

func BatchAssignMembersToCluster(ctx context.Context, db *gorm.DB, memberIds []int32, clusterId int32) error {
	mp := map[string]interface{}{
		"cluster_id": clusterId,
	}
	// FIXME(@yfaming): 此处会生成错误的 SQL 语句，形如：
	// UPDATE `dbcluster_member` SET `cluster_id` = 1  WHERE (id IN 1,2)
	// 接下来考虑升级 gorm 解决。目前所依赖为 github.com/jinzhu/gorm，应升级到 gorm.io/gorm
	return db.Debug().Table("dbcluster_member").Where("id IN ?", memberIds).Updates(mp).Error
}

func BatchUnassignClusterMembers(ctx context.Context, db *gorm.DB, memberIds []int32) error {
	return BatchAssignMembersToCluster(ctx, db, memberIds, 0)
}
