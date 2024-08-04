package websocket

import (
	"fmt"
	"mygowebsockt/helper"
	"mygowebsockt/models"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	defaultAppID = 101 //默认应用ID
)

var (
	clientManager = NewClientManager()
	appIDs        = []uint32{defaultAppID, 102, 103, 104}
	serverIp      string
	serverPort    string
)

func GetAppIDs() []uint32 {
	return appIDs
}

func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)
	return
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}
	return
}

func InAppIDs(appID uint32) (inAppIDs bool) {
	for _, value := range appIDs {
		if value == appID {
			inAppIDs = true
			return
		}
	}
	return
}

func GetDefaultAppID() (appID uint32) {
	appID = defaultAppID
	return
}

func StartWebSocket() {
	serverIp = helper.GetServerIp()
	websocketPort := viper.GetString("app.websocketPort")
	rpcPort := viper.GetString("app.rpcPort")
	serverPort = rpcPort
	http.HandleFunc("/acc", wsPage)
	go clientManager.start()
	fmt.Println("wevSocket 启动程序成功", serverIp, serverPort)
	_ = http.ListenAndServe(":"+websocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])
		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)
	go client.read()
	go client.write()

	clientManager.Register <- client
}
