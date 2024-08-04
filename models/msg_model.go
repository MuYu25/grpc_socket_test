package models

import "mygowebsockt/setting"

const (
	// MessageTypeText 消息类型：文本
	MessageTypeText = "text"
	// MessageCmdMsg 文本类型消息
	MessageCmdMsg = "msg"
	// MessageCmdEnter 进入聊天室
	MessageCmdEnter = "enter"
	// MessageCmdExit 退出聊天室
	MessageCmdExit = "exit"
)

type Message struct {
	Target string `json:"target"` // 目标用户
	Type   string `json:"type"`   // 消息类型 text/img
	Msg    string `json:"msg"`    // 消息内容
	From   string `json:"from"`   // 发送者
}

func NesMsg(from string, Msg string) *Message {
	return &Message{
		Type: MessageTypeText,
		Msg:  Msg,
		From: from,
	}
}

func getTextMsgData(cmd, uuID, msgID, message string) string {
	textMsg := NesMsg(uuID, message)
	head := NewResponseHead(msgID, cmd, setting.OK, "OK", textMsg)
	return head.String()
}

func GetMsgData(uuID, msgID, cmd, message string) string {
	return getTextMsgData(cmd, uuID, msgID, message)
}

func GetTextMsgData(uuID, msgID, message string) string {
	return getTextMsgData("msg", uuID, msgID, message)
}

func GetTextMsgDataEnter(uuID, msgID, message string) string {
	return getTextMsgData("enter", uuID, msgID, message)
}

func GetTextMsgDataExit(uuID, msgID, message string) string {
	return getTextMsgData("exit", uuID, msgID, message)
}
