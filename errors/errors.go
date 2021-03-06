package errors

import (
	"encoding/json"
	"fmt"
)

const (
	// 业务错误码从 MinCustomErrorCode 开始，与 HTTP status code 隔开。
	// Code < MinCustomErrorCode : http 返回错误码 Code
	// Code >= MinCustomErrorCode: http 返回错误码 http.StatusBadRequest(400)
	MinCustomErrorCode = 1000
)

type (
	ErrInfo struct {
		Code int32  `json:"code"`
		Name string `json:"name"`
		Msg  string `json:"msg"`
	}
)

func newErrInfo(code int32, name string, msg string) *ErrInfo {
	return &ErrInfo{
		Code: code,
		Name: name,
		Msg:  msg,
	}
}

func (err *ErrInfo) Error() string {
	res, _ := json.Marshal(err)
	return string(res)
}

// 1000-1999 filebeat 相关错误码
func FileBeatConfError(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(1000, "FilebeatConfError", fmt.Sprintf(format, a...))
}

func FilebeatRateLimited(msg string) *ErrInfo {
	return newErrInfo(1001, "FilebeatRateLimited", msg)
}

func FilebeatSearchParamsError(msg string) *ErrInfo {
	return newErrInfo(1002, "FilebeatSearchParamsError", msg)
}

func FilebeatSearchError(msg string) *ErrInfo {
	return newErrInfo(1003, "FilebeatSearchError", msg)
}

// 2000 - 2999 dbcluster 相关
func InvalidDbType(msg string) *ErrInfo {
	return newErrInfo(2000, "InvalidDbType", msg)
}

func DbClusterNotFound(msg string) *ErrInfo {
	return newErrInfo(2001, "DbClusterNotFound", msg)
}
func DbClusterNotFoundById(id int32) *ErrInfo {
	msg := fmt.Sprintf("dbcluster with id=%d not found", id)
	return DbClusterNotFound(msg)
}
func DbClusterNotFoundByName(name string) *ErrInfo {
	msg := fmt.Sprintf("dbcluster with name=%s not found", name)
	return DbClusterNotFound(msg)
}

// 3000-3999 agent 相关错误
func AgentConfigError(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(3000, "AgentConfError", fmt.Sprintf(format, a...))
}

func AgentReportError(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(3001, "AgentReportError", fmt.Sprintf(format, a...))
}

func AgentHeartbeatError(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(3002, "AgentHeartbeatError", fmt.Sprintf(format, a...))
}

// 4000-4999 dbcluster member 相关
const DBClusterMemberNotFoundErrorCode = 4000

func IsDBClusterMemberNotFound(err error) bool {
	errInfo, ok := err.(*ErrInfo)
	if !ok {
		return false
	}
	return errInfo.Code == DBClusterMemberNotFoundErrorCode
}
func DbClusterMemberNotFound(msg string) *ErrInfo {
	return newErrInfo(DBClusterMemberNotFoundErrorCode, "DbClusterMemberNotFound", msg)
}
func DbClusterMemberNotFoundById(id int32) *ErrInfo {
	msg := fmt.Sprintf("dbcluster_member with id=%d not found", id)
	return DbClusterMemberNotFound(msg)
}
func DbClusterMemberNotFoundByIpAndPort(ipAddr string, port int16) *ErrInfo {
	msg := fmt.Sprintf("dbcluster_member with ip=%s and port=%d not found", ipAddr, port)
	return DbClusterMemberNotFound(msg)
}

func DbClusterMemberNotAssigned(msg string) *ErrInfo {
	return newErrInfo(4001, "DbClusterMemberNotAssigned", msg)
}

// 5000-5999 tag 相关
func InvalidItemType(msg string) *ErrInfo {
	return newErrInfo(5000, "InvalidItemType", msg)
}

// 6000-6999 es index 相关
func InvalidTimeDuration(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(6000, "InvalidTimeDuration", fmt.Sprintf(format, a...))
}

func ConnectToESFailed(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(6001, "ConnectToESFailed", fmt.Sprintf(format, a...))
}
