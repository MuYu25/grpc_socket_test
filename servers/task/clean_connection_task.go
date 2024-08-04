package task

import (
	"fmt"
	"mygowebsockt/servers/websocket"
	"runtime/debug"
	"time"
)

func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)
}

func cleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnection stop", r, string(debug.Stack()))
		}
	}()
	fmt.Println("定时任务，清理超时连接", param)
	websocket.ClearTimeoutConnections()
	return
}
