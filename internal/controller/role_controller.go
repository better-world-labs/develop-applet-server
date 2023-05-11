package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

type Base struct {
	gone.Flag
	PlanetsService service.IPlanet `gone:"*"`
}

//go:gone
func NewBaseController() gone.Goner {
	return &Base{}
}

func (ctr Base) WithRole(role entity.PlanetRole, f func(ctx *gin.Context) (any, error)) func(ctx *gin.Context) (any, error) {
	return func(ctx *gin.Context) (any, error) {
		currentPlanetId := 1 // 默认为 1, 后续会让前端传递
		userId := utils.CtxMustGetUserId(ctx)

		r, err := ctr.PlanetsService.GetPlanetRoles(int64(currentPlanetId), userId)
		if err != nil {
			return nil, gin.NewBusinessError(err.Error(), 403)
		}

		if r < role {
			return nil, gin.NewBusinessError("permission denied", 403)
		}

		return f(ctx)
	}
}
