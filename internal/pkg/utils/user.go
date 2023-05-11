package utils

import (
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

// CtxMustGetUserId 获取 gin.Context 中 存储的 CreatedBy,不存在则 panic
func CtxMustGetUserId(ctx *gin.Context) int64 {
	value, exists := ctx.Get(entity.UserIdKey)
	if !exists || value == nil {
		panic("userId not found in context")
	}

	return value.(int64)
}
