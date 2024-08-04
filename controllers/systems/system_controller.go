package systems

import (
	"fmt"
	"mygowebsockt/controllers"
	"mygowebsockt/servers/websocket"
	"mygowebsockt/setting"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Staus(c *gin.Context) {
	isDebug := c.Query("debug")
	fmt.Println("http_request 查询系统状态", isDebug)
	data := make(map[string]interface{})
	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager 的数量
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)
	controllers.Response(c, setting.OK, "", data)
}
