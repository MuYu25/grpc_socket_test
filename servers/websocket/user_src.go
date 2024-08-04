package websocket

import (
	"errors"
	"fmt"
	"mygowebsockt/lib/cache"
	"mygowebsockt/models"
	"mygowebsockt/servers/grpcclient"
	"time"

	"github.com/redis/go-redis/v9"
)

func UserList(appID uint32) (userList []string) {
	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发送消息", err)
		return
	}
	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(appID)
		} else {
			list, _ = grpcclient.GetUserList(server, appID)
		}
		userList = append(userList, list...)
	}
	return
}

func CheckUserOnline(appID uint32, userID string) (online bool) {
	if appID == 0 {
		for _, appID := range GetAppIDs() {
			online, _ = checkUserOnline(appID, userID)
			if online {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appID, userID)
	}
	return
}

func checkUserOnline(appID uint32, userID string) (online bool, err error) {
	key := GetUserKey(appID, userID)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("GetUserOnlineInfo", appID, userID, err)
			return false, nil
		}
		fmt.Println("GetUserOnlineInfo", appID, userID, err)
		return
	}
	online = userOnline.IsOnline()
	return
}

func SendUserMessage(appID uint32, userID string, msgID, message string) (sendResult bool, err error) {
	data := models.GetTextMsgData(userID, msgID, message)
	client := GetUserClient(appID, userID)
	if client != nil {
		sendResult, err = SendUserMessageLocal(appID, userID, data)
		if err != nil {
			fmt.Println("给用户发送消息", appID, userID, err)
		}
		return
	}
	key := GetUserKey(appID, userID)
	info, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		fmt.Println("给用户发送消息fail", appID, userID, err)
		return false, nil
	}
	if !info.IsOnline() {
		fmt.Println("用户不在线", key)
		return false, nil
	}
	server := models.NewServer(info.AccIp, info.AccPort)
	msg, err := grpcclient.SendMsg(server, msgID, appID, userID, models.MessageCmdMsg, models.MessageCmdMsg, message)
	if err != nil {
		fmt.Println("给用户发送消息成功fail", key, err)
		return false, err
	}
	fmt.Println("给用户发送消息成功-rpc", msg)
	sendResult = true
	return
}

func SendUserMessageLocal(appID uint32, userID string, data string) (sendResult bool, err error) {
	client := GetUserClient(appID, userID)
	if client == nil {
		err = fmt.Errorf("用户不在线")
		return
	}
	client.SendMsg([]byte(data))
	sendResult = true
	return
}

func SendUserMessageAll(appID uint32, userID string, msgID, cmd, message string) (sendResult bool, err error) {
	sendResult = true
	currrentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currrentTime)
	if err != nil {
		fmt.Println("给全体用户发送消息", err)
		return
	}
	for _, server := range servers {
		if IsLocal(server) {
			data := models.GetMsgData(userID, msgID, cmd, message)
			AllSendMessages(appID, userID, data)
		} else {
			_, _ = grpcclient.SendMsgAll(server, msgID, appID, userID, cmd, message)
		}
	}
	return
}
