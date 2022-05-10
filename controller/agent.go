package controller

import (
	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/model/es"
	"github.com/dbrainhub/dbrainhub/server"
	"github.com/dbrainhub/dbrainhub/utils/logger"
	"github.com/gin-gonic/gin"
)

func Heartbeat(c *gin.Context, req *api.HeartbeatRequest) (*api.HeartbeatResponse, error) {
	logger.Infof("receive heartbeat from IP:%s, port: %d, req: %s\n", req.AgentInfo.Localip, req.DbInfo.Port, req.String())

	if req.AgentInfo.Localip == "" || req.DbInfo.Port == 0 {
		logger.Errorf("invalid agent ip.port:%s.%d", req.AgentInfo.Localip, req.DbInfo.Port)
		return nil, errors.AgentHeartbeatError("invalid ip.port")
	}

	member, err := model.GetDbClusterMemberByIpAndPort(c, model.GetDB(c), req.AgentInfo.Localip, int16(req.DbInfo.Port))
	if err != nil && errors.IsDBClusterMemberNotFound(err) {
		logger.Errorf("unknown agent heartbeat ip.port: %s.%d", req.AgentInfo.Localip, req.DbInfo.Port)
		return nil, errors.AgentHeartbeatError("unknown agent")
	}
	if err != nil {
		logger.Errorf("GetDbClusterMemberByIpAndPort error when heartbeat, err: %v, req: %#v", err, req)
		return nil, errors.AgentHeartbeatError("visit db error")
	}

	cluster, err := model.GetDbClusterById(c, model.GetDB(c), member.ClusterId)
	if err != nil {
		logger.Errorf("GetDbClusterById error when find cluster, err: %v, req: %#v", err, req)
		return nil, err
	}

	port := int16(req.DbInfo.Port)
	updateParam := &model.UpdateDbClusterMemberParams{
		Id:     member.Id,
		IPaddr: &req.AgentInfo.Localip,
		Port:   &port,
	}
	if err := model.UpdateDbClusterMember(c, model.GetDB(c), updateParam); err != nil {
		logger.Errorf("UpdateDbClusterMember error when heartbeat, err: %v, req: %#v", err, req)
		return nil, errors.AgentHeartbeatError("update cluster member failed")
	}

	server.GetDefaultAsyncESClient().Send([]*es.ESMessage{{
		Meta: &es.ESMeta{
			Index: server.GetIndicesIndexName(),
		},
		Data: &es.AgentIndexData{
			TimeStamp: req.AgentInfo.Datetime,
			IP:        req.AgentInfo.Localip,
			Port:      int(req.DbInfo.Port),
			CPURatio:  req.AgentInfo.CpuRatio,
			MemRatio:  req.AgentInfo.MemRatio,
			DiskRatio: req.AgentInfo.DiskRatio,
			QPS:       req.DbInfo.Qps,
			TPS:       req.DbInfo.Tps,
			Cluster:   cluster.Name,
		},
	}})

	return &api.HeartbeatResponse{}, nil
}

func Report(c *gin.Context, req *api.StartupReportRequest) (*api.StartupReportResponse, error) {
	logger.Infof("receive reporter from IP:%s, port: %d, req: %s\n", req.IpAddr, req.Port, req.String())

	if req.IpAddr == "" || req.Port == 0 {
		logger.Errorf("invalid agent ip.port:%s.%d", req.IpAddr, req.Port)
		return nil, errors.AgentReportError("invalid ip.port")
	}
	member, err := model.GetDbClusterMemberByIpAndPort(c, model.GetDB(c), req.IpAddr, int16(req.Port))
	if err != nil && !errors.IsDBClusterMemberNotFound(err) {
		logger.Errorf("GetDbClusterMemberByIpAndPort error when report, err: %v, req: %#v", err, req)
		return nil, errors.AgentReportError("visit db error")
	}

	insertParam := &model.NewDbClusterMember{
		ClusterId: 0,
		Hostname:  req.Hostname,
		DbType:    req.DbType.String(),
		DbVersion: req.DbVersion,
		Role:      0,
		IPaddr:    req.IpAddr,
		Port:      int16(req.Port),
		Os:        req.Os,
		OsVersion: req.OsVersion,
		HostType:  int32(req.HostType),
		Env:       req.Env,
	}
	if member == nil {
		if _, err := model.CreateDbClusterMember(c, model.GetDB(c), insertParam); err != nil {
			logger.Errorf("CreateDbClusterMember error when report, err: %v, req: %#v", err, req)
			return nil, errors.AgentReportError("insert new cluster member failed")
		}
	} else {
		updateParam := &model.UpdateDbClusterMemberParams{
			Id:        member.Id,
			Hostname:  &insertParam.Hostname,
			DbType:    &insertParam.DbType,
			DbVersion: &insertParam.DbVersion,
			Role:      &insertParam.Role,
			IPaddr:    &insertParam.IPaddr,
			Port:      &insertParam.Port,
			Os:        &insertParam.Os,
			OsVersion: &insertParam.OsVersion,
			HostType:  &insertParam.HostType,
			Env:       &insertParam.Env,
		}
		if err := model.UpdateDbClusterMember(c, model.GetDB(c), updateParam); err != nil {
			logger.Errorf("UpdateDbClusterMember error when report, err: %v, req: %#v", err, req)
			return nil, errors.AgentReportError("update cluster member failed")
		}
	}
	return &api.StartupReportResponse{}, nil
}
