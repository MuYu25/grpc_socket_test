package models

type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一ID
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` //数据json
}

type Login struct {
	ServiceToken string `json:"serviceToken"`
	AppID        uint32 `json:"appID,omitempty"`
	UserID       string `json:"userID,omitempty"`
}

type HeartBeat struct {
	UserID string `json:"userID,omitempty"`
}
