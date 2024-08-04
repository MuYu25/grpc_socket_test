package user

import (
	"fmt"
	"mygowebsockt/controllers"
	"mygowebsockt/lib/cache"
	"mygowebsockt/models"
	"mygowebsockt/servers/websocket"
	"mygowebsockt/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseUint(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看全部在线用户", appID)
	data := make(map[string]interface{})
	userList := websocket.UserList(appID)
	data["userList"] = userList
	data["userCount"] = len(userList)
	controllers.Response(c, setting.OK, "", data)
}

func Online(c *gin.Context) {
	userID := c.Query("userID")
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseUint(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 查看用户是否在线", userID, appIDStr)
	data := make(map[string]interface{})
	online := websocket.CheckUserOnline(appID, userID)
	data["userID"] = userID
	data["online"] = online
	controllers.Response(c, setting.OK, "", data)
}

func SendMessage(c *gin.Context) {
	appIDStr := c.Query("appID")
	userID := c.Query("userID")
	msgID := c.Query("msgID")
	message := c.PostForm("message")
	appIDUint64, _ := strconv.ParseUint(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 发送消息", userID, appIDStr, msgID, message)
	// 暂时未使用token进行验证
	data := make(map[string]interface{})
	if cache.SeqDupliCates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, setting.OK, "重复提交", data)
		return
	}
	sendResult, err := websocket.SendUserMessage(appID, userID, msgID, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResult
	controllers.Response(c, setting.OK, "", data)
}

func SendMessageAll(c *gin.Context) {
	appIDStr := c.PostForm("appID")
	userID := c.PostForm("userID")
	msgID := c.PostForm("msgID")
	message := c.PostForm("message")
	appIDUint64, _ := strconv.ParseUint(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	fmt.Println("http_request 发送消息", userID, appIDStr, msgID, message)
	data := make(map[string]interface{})
	if cache.SeqDupliCates(msgID) {
		fmt.Println("给用户发送消息 重复提交:", msgID)
		controllers.Response(c, setting.OK, "重复提交", data)
		return
	}
	sendResult, err := websocket.SendUserMessageAll(appID, userID, msgID, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}
	data["sendResults"] = sendResult
	controllers.Response(c, setting.OK, "", data)
}
