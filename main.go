package main

import (
	"fmt"
	"io"
	"mygowebsockt/lib/redislib"
	"mygowebsockt/routers"
	"mygowebsockt/servers/grpcserver"
	"mygowebsockt/servers/task"
	"mygowebsockt/servers/websocket"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	initFile()
	initRedis()
	router := gin.Default()
	routers.Init(router)
	routers.WebSocketInit()

	task.Init()

	task.ServerInit()
	go websocket.StartWebSocket()

	go grpcserver.Init()
	go open()
	httpPort := viper.GetString("app.httpPort")
	_ = http.ListenAndServe(":"+httpPort, router)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")    // 指定查找配置文件的路径
	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))
}

func initFile() {
	gin.DisableConsoleColor()
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initRedis() {
	redislib.NewClinet()
}

func open() {
	time.Sleep(1000 * time.Millisecond)
	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"
	fmt.Println("访问页面体验", httpUrl)
	cmd := exec.Command("open", httpUrl)
	_, _ = cmd.Output()
}
