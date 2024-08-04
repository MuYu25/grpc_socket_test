package models

import "encoding/json"

type Head struct {
	Seq      string    `json:"seq`
	Cmd      string    `json:"cmd"`
	Response *Response `json:"response"`
}

type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"`
}

type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)
	return &Head{
		Seq:      seq,
		Cmd:      cmd,
		Response: response,
	}
}

func (h *Head) String() (headstr string) {
	headBytes, _ := json.Marshal(h)
	headstr = string(headBytes)
	return
}

func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{
		Code:    code,
		CodeMsg: codeMsg,
		Data:    data,
	}
}
