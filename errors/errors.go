package errors

import (
	"encoding/json"
)

type (
	// 业务错误码从 1000 开始，与 HTTP status code 隔开
	ErrInfo struct {
		Code int32  `json:"code"`
		Name string `json:"name"`
		Msg  string `json:"msg"`
	}

	SuccInfo struct {
		Code int32       `json:"code"`
		Data interface{} `json:"data"`
	}
)

func NewUnknownError() error {
	return NewErrInfo(-1, "internal-error", "unknown error")
}

func NewSuccessResp(v interface{}) *SuccInfo {
	return &SuccInfo{
		Code: 0,
		Data: v,
	}
}

func NewErrInfo(code int32, name string, msg string) *ErrInfo {
	return &ErrInfo{
		Code: code,
		Name: name,
		Msg:  msg,
	}
}

func (s *SuccInfo) String() string {
	res, _ := json.Marshal(s)
	return string(res)
}

func (err *ErrInfo) Error() string {
	res, _ := json.Marshal(err)
	return string(res)
}
