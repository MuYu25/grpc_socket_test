package routers

import (
	"mygowebsockt/controllers/home"
	"mygowebsockt/controllers/systems"
	"mygowebsockt/controllers/user"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.LoadHTMLGlob("views/**/*")

	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", user.List)
		userRouter.GET("/online", user.Online)
		userRouter.POST("/sendMessage", user.SendMessage)
		userRouter.POST("/sendMessageAll", user.SendMessageAll)
	}

	systemRouter := router.Group("/system")
	{
		systemRouter.GET("/state", systems.Staus)
	}

	homeRouter := router.Group("/home")
	{
		homeRouter.GET("/index", home.Index)
	}
}
