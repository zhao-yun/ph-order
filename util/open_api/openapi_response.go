package open_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandResponse struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// OpenApiSuccessResponse 处理接口正确返回.
func OpenApiSuccessResponse(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}

	res := StandResponse{
		Code: 0,
		Msg:  "ApiSuccessResponse",
		Data: data,
	}

	ctx.JSON(http.StatusOK, res)
	return

}

// OpenApiErrorResponse 处理接口错误返回.
func OpenApiErrorResponse(ctx *gin.Context, code int, message string) {

	res := StandResponse{
		Code: int64(code),
		Msg:  message,
	}

	ctx.JSON(code, res)
	return
}
