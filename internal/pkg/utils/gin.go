package utils

import (
	"github.com/gone-io/gone/goner/gin"
	"strconv"
)

// CtxPathParamInt64 获取 gin.Context 中 int64 类型 PathParam
func CtxPathParamInt64(ctx *gin.Context, key string) (int64, error) {
	param := ctx.Param(key)
	return strconv.ParseInt(param, 10, 64)
}

// CtxGetString 获取 gin.Context 中 存储的 String 值
func CtxGetString(context *gin.Context, key string) string {
	s, existed := context.Get(key)
	if !existed {
		return ""
	}
	return s.(string)
}
