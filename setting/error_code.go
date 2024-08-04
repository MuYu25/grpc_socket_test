package setting

const (
	OK                 = 200  //Success
	NotLoggedIn        = 1000 // 未登录
	ParameterIllegal   = 1001 // 参数不合法
	UnauthorizedUserID = 1002 // 用户ID不合法
	Unauthorized       = 1003 // 未授权
	ServerError        = 1004 // 服务器错误
	NotData            = 1005 // 没有数据
	NodelAddError      = 1006 // 添加错误
	ModelDeleteError   = 1007 // 删除错误
	ModelStoreError    = 1008 // 存储错误
	OperationFailure   = 1009 // 操作失败
	RoutingNoExist     = 1010 // 路由不存在
)

func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		OK:                 "Success",
		NotLoggedIn:        "未登录",
		ParameterIllegal:   " 参数不合法",
		UnauthorizedUserID: "用户ID不合法",
		Unauthorized:       "未授权",
		ServerError:        " 服务器错误",
		NotData:            " 没有数据",
		NodelAddError:      " 添加错误",
		ModelDeleteError:   " 删除错误",
		ModelStoreError:    " 存储错误",
		OperationFailure:   " 操作失败",
		RoutingNoExist:     " 路由不存在",
	}
	if message == "" {
		if value, ok := codeMap[code]; ok {
			codeMessage = value
		} else {
			codeMessage = "未定义错误类型!"
		}
	} else {
		codeMessage = message
	}
	return codeMessage
}
