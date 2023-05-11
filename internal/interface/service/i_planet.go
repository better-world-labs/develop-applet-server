package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IPlanet interface {
	UpdatePlanetMessage(icon entity.Url, frontCover entity.Url, name string, planetId int64) error

	PagePlanetMembers(planetId int64, query page.Query, userIdLikeCondition string) (page.Result[*domain.Member], error)

	CountPlanetMembers(planetId int64) (int64, error)

	ListPlanetAdmins(planetId int64) ([]*entity.PlanetMember, error)

	UpdateMembersRole(planetId int64, userIds []int64, role entity.PlanetRole) error

	MemberJoin(userId, planetId int64) error

	UpdateMembersStatus(planetId int64, userIds []int64, status entity.PlanetMemberStatus) error

	GetPlanetMemberByUserId(planetId, userId int64) (*entity.PlanetMember, error)

	GetPlanetRoles(planetId, userId int64) (entity.PlanetRole, error)

	ListPlanetRolesMap(planetId int64, userIds []int64) (map[int64]entity.PlanetRole, error)

	GetPlanet(planetId int64) (*entity.Planet, error)
}
