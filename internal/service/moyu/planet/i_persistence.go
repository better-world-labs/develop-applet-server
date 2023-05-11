package planet

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type iPersistence interface {
	updatePlanetMessageByPlanetId(planet entity.Planet) error

	getPlanetById(planetId int64) (*entity.Planet, error)

	pageMembersByPlanetId(planetId int64, query page.Query, userIdLikeCondition string) (page.Result[*entity.PlanetMember], error)

	batchUpdateMembersRole(planetId int64, userIds []int64, role entity.PlanetRole) error

	batchUpdateMembersStatus(planetId int64, userIds []int64, role entity.PlanetMemberStatus) error

	getMember(planetId, userId int64) (*entity.PlanetMember, error)

	createMember(member *entity.PlanetMember) error

	listAdmins(planetId int64) ([]*entity.PlanetMember, error)

	listPlanetMembersByUserIds(planetId int64, userIds []int64) ([]*entity.PlanetMember, error)

	countMembersByPlanetId(planetId int64) (int64, error)
}
