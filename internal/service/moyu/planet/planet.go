package planet

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"time"
)

//go:gone
func NewPlanetService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Flag
	P           iPersistence  `gone:"*"`
	UserService service.IUser `gone:"*"`
}

func (s *svc) UpdatePlanetMessage(icon entity.Url, frontCover entity.Url, name string, planetId int64) error {
	plant := entity.Planet{
		Id:         planetId,
		Name:       name,
		Icon:       string(icon),
		FrontCover: string(frontCover),
	}
	return s.P.updatePlanetMessageByPlanetId(plant)
}

func (s *svc) CountPlanetMembers(planetId int64) (int64, error) {
	return s.P.countMembersByPlanetId(planetId)
}

func (s *svc) PagePlanetMembers(planetId int64, query page.Query, userIdLikeCondition string) (page.Result[*domain.Member], error) {
	var result page.Result[*domain.Member]

	members, err := s.P.pageMembersByPlanetId(planetId, query, userIdLikeCondition)
	if err != nil {
		return result, err
	}

	result.Total = members.Total
	for _, aMember := range members.List {
		aUser, err := s.UserService.GetUserById(aMember.UserId)
		if err != nil {
			return result, err
		}

		m := domain.Member{
			User: domain.User{
				Id:       aUser.Id,
				Nickname: aUser.Nickname,
				Avatar:   string(aUser.Avatar),
				Online:   aUser.Online,
			},
			Role:   int64(aMember.Role),
			Status: int64(aMember.Status),
		}
		result.List = append(result.List, &m)
	}

	return result, nil
}

func (s *svc) ListPlanetAdmins(planetId int64) ([]*entity.PlanetMember, error) {
	return s.P.listAdmins(planetId)
}

func (s *svc) UpdateMembersRole(planetId int64, userIds []int64, role entity.PlanetRole) error {
	return s.P.batchUpdateMembersRole(planetId, userIds, role)
}

func (s *svc) UpdateMembersStatus(planetId int64, userIds []int64, status entity.PlanetMemberStatus) error {
	return s.P.batchUpdateMembersStatus(planetId, userIds, status)
}

func (s *svc) ListPlanetRolesMap(planetId int64, userIds []int64) (map[int64]entity.PlanetRole, error) {
	members, err := s.P.listPlanetMembersByUserIds(planetId, userIds)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(members, func(member *entity.PlanetMember) (int64, entity.PlanetRole) {
		return member.UserId, member.Role
	}), err
}

func (s *svc) GetPlanetRoles(planetId, userId int64) (entity.PlanetRole, error) {
	memberInfo, err := s.GetPlanetMemberByUserId(planetId, userId)
	if err != nil {
		return 0, err
	}

	return memberInfo.Role, nil
}

func (s *svc) GetPlanetMemberByUserId(planetId, userId int64) (*entity.PlanetMember, error) {
	return s.P.getMember(planetId, userId)
}

func (s *svc) MemberJoin(userId, planetId int64) error {
	u, err := s.UserService.GetUserById(userId)
	if err != nil {
		return err
	}

	if u == nil {
		return gin.NewParameterError("user not found")
	}

	p, err := s.P.getPlanetById(planetId)
	if err != nil {
		return err
	}

	if p == nil {
		return gin.NewParameterError("planet not found")
	}

	return s.addMember(userId, planetId)
}

func (s *svc) addMember(userId, planetId int64) error {
	member, err := s.P.getMember(planetId, userId)
	if err != nil {
		return err
	}

	if member != nil {
		return gin.NewBusinessError("user has joined this planet")
	}

	return s.P.createMember(&entity.PlanetMember{
		PlanetId:  planetId,
		UserId:    userId,
		Role:      entity.PlanetRoleMember,
		Status:    entity.PlanetMemberStatusOK,
		CreatedAt: time.Now(),
	})

}

func (s *svc) GetPlanet(planetId int64) (*entity.Planet, error) {
	return s.P.getPlanetById(planetId)
}
