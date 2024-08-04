package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"mygowebsockt/lib/cache"
	"mygowebsockt/models"
	"mygowebsockt/setting"
	"time"

	"github.com/redis/go-redis/v9"
)

func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = setting.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)
	data = "pong"
	return
}

func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = setting.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = setting.ParameterIllegal
		fmt.Println("用户登陆 非法的用户", seq, request.UserID)
		return
	}

	// 用户权限验证， 一般是token
	if request.UserID == "" || len(request.UserID) >= 20 {
		code = setting.UnauthorizedUserID
		fmt.Println("用户登陆 非法的用户", seq, request.UserID)
		return
	}
	if !InAppIDs(request.AppID) {
		code = setting.Unauthorized
		fmt.Println("用户登陆 不支持的平台", client.AppID, client.UserID, seq)
		return
	}
	if client.IsLogin() {
		fmt.Println("用户登陆 用户已经登陆", client.AppID, client.UserID, seq)
		code = setting.OperationFailure
		return
	}
	client.Login(request.AppID, request.UserID, currentTime)

	userOnlie := models.UserLogin(serverIp, serverPort, request.AppID, request.UserID, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnlie)
	if err != nil {
		code = setting.ServerError
		fmt.Println("用户登陆 SetUserOnlineInfo", seq, err)
		return
	}

	login := &login{
		AppID:  request.AppID,
		UserID: request.UserID,
		Client: client,
	}
	clientManager.Login <- login
	fmt.Println("用户登陆 成功", seq, client.Addr, request.UserID)
	return
}

func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = setting.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = setting.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)
		return
	}
	fmt.Println("webSocket_request 心跳接口", client.AppID, client.UserID)
	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登陆", client.AppID, client.UserID, seq)
		code = setting.NotLoggedIn
		return
	}
	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			code = setting.NotLoggedIn
			fmt.Println("心跳接口 GetUserOnlineInfo", seq, client.AppID, client.UserID, err)
			return
		}
	}
	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = setting.ServerError
		fmt.Println("心跳接口 SetUserOnlineInfo", seq, client.AppID, client.UserID, err)
		return
	}
	return
}
