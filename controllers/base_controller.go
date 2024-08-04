package controllers

import (
	"mygowebsockt/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	gin.Context
}

func Response(c *gin.Context, code uint32, msg string, data map[string]interface{}) {
	message := setting.Response(code, msg, data)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, session")
	c.Set("content-type", "application/json")
	c.JSON(http.StatusOK, message)
}
