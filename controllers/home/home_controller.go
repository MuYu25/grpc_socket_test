package home

import (
	"fmt"
	"mygowebsockt/servers/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Index(c *gin.Context) {
	appIDStr := c.Query("appID")
	appIDUint64, _ := strconv.ParseUint(appIDStr, 10, 32)
	appID := uint32(appIDUint64)
	if !websocket.InAppIDs(appID) {
		appID = websocket.GetDefaultAppID()
	}
	fmt.Println("http_request 聊天首页", appID)
	data := gin.H{
		"title":        "聊天首页",
		"appID":        appID,
		"httpUrl":      viper.GetString("app.httpUrl"),
		"webSocketUrl": viper.GetString("app.websocketUrl"),
	}
	c.HTML(http.StatusOK, "index.tpl", data)
}
