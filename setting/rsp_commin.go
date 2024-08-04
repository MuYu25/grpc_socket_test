package setting

import "github.com/spf13/viper"

type JSONResult struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(code uint32, messgage string, data interface{}) JSONResult {
	messgage = GetErrorMessage(code, messgage)
	jsonMap := grantMap(code, messgage, data)
	return jsonMap
}

func grantMap(code uint32, message string, data interface{}) JSONResult {
	return JSONResult{
		Code: code,
		Msg:  message,
		Data: data,
	}
}

func Seting() map[string]interface{} {
	m := map[string]interface{}{
		"message": viper.GetString("app.logFile"),
	}
	return m
}
