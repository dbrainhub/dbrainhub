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

// 1000-2000 filebeat 相关错误码
func FileBeatConfError(format string, a ...interface{}) *ErrInfo {
	return newErrInfo(1000, "filebeat-conf-error", fmt.Sprintf(format, a...))
}
