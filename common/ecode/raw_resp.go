package ecode

import "encoding/json"

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponseMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func (resp *RespMsg) JOSNBytes() ([]byte, error) {
	r, err := json.Marshal(resp)
	if err != nil {
		return []byte(""), err
	}
	return r, nil
}

func (resp *RespMsg) JOSNString() (string, error) {
	r, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
