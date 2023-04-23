package scikits

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSystemErr        = -1  // 系统错误，需要修复bug
	CodeOK               = 0   // 成功
	CodeShowErr          = 2   // 前端需要弹窗展示的错误
	CodeLoginExpired     = 10  // 登录过期
	CodePermissionError  = 20  // 没有权限
	CodeJumpBoundLibrary = 23  // 跳转绑定图书馆用户页面
	CodeRedirectSciBrain = 301 // 重新定向页面
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
	CodeLoginExpired: ResponseFormat{
		Code:    CodeLoginExpired,
		Message: "Login Expire",
	},
	CodePermissionError: ResponseFormat{
		Code:    CodePermissionError,
		Message: "Permission Limited",
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
