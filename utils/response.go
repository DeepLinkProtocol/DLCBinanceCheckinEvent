package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TaskCompletionData struct {
	Sign bool `json:"sign"`
}

type SignCountData struct {
	Count int `json:"count"`
}

func JSONResponse(c *gin.Context, statusCode int, code, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
