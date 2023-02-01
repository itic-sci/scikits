package scikits

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSystemErr        = -1 // 系统错误，需要修复bug
	CodeOK               = 0  // 成功
	CodeParamErr         = 1  // 传的参数错误，不在正常范围内
	CodeShowErr          = 2  // 前端需要弹窗展示的错误
	CodeApiCallLimited   = 10 // 权限错误
	CodeLoginExpire      = 22 // 登录过期
	CodeJumpBoundLibrary = 23 // 跳转绑定图书馆用户页面
)

type ResponseFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var codeMap = map[int]ResponseFormat{
	CodeOK: ResponseFormat{
		Code:    CodeOK,
		Message: "OK",
	},
	CodeSystemErr: ResponseFormat{
		Code:    CodeSystemErr,
		Message: "system error",
	},
	CodeParamErr: ResponseFormat{
		Code:    CodeParamErr,
		Message: "Param error",
	},
	CodeApiCallLimited: ResponseFormat{
		Code:    CodeApiCallLimited,
		Message: "Api Call Limited",
	},
	CodeLoginExpire: ResponseFormat{
		Code:    CodeLoginExpire,
		Message: "Login Expire",
	},
}

func RespondError(c *gin.Context, code int, errMsgArr ...string) {
	errMsg := ""
	if len(errMsgArr) > 0 {
		errMsg = errMsgArr[0]
	}
	ec, ok := codeMap[code]
	var message string
	if ok {
		message = ec.Message + " " + errMsg
	} else {
		message = errMsg
	}

	r := ResponseFormat{
		Code:    code,
		Message: message,
	}

	c.JSON(http.StatusOK, r)
}

func RespondFormatData(c *gin.Context, data interface{}) {
	r := ResponseFormat{
		Code:    CodeOK,
		Message: "Ok",
		Data:    data,
	}
	c.JSON(http.StatusOK, r)
}

func RespondData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
