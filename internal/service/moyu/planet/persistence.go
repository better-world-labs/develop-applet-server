package planet

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/jmoiron/sqlx"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

const (
	RootTableName   = "planet"
	MemberTableName = "planet_member"
)

//go:gone
func NewPlanetPersistence() gone.Goner {
	return &persistence{}
}

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *persistence) updatePlanetMessageByPlanetId(planet entity.Planet) error {
	_, err := p.Table(RootTableName).ID(planet.Id).Cols("icon", "front_cover", "name").Update(planet)
	return err
}

func (p *persistence) countMembersByPlanetId(planetId int64) (int64, error) {
	return p.Table(MemberTableName).Where("planet_id = ?", planetId).Count()
}

func (p *persistence) pageMembersByPlanetId(planetId int64, query page.Query, userIdLikeCondition string) (page.Result[*entity.PlanetMember], error) {
	var result page.Result[*entity.PlanetMember]

	session := p.Table(MemberTableName).Where("planet_id = ?", planetId)
	if userIdLikeCondition != "" {
		session.And("user_id like concat('%',?,'%')", userIdLikeCondition)
	}

	count, err := session.Limit(query.LimitOffset(), query.LimitStart()).FindAndCount(&result.List)
	if err != nil {
		return result, err
	}

	result.Total = count
	return result, err
}

func (p *persistence) batchUpdateMembersRole(planetId int64, userIds []int64, role entity.PlanetRole) error {
	sql, args, err := sqlx.In(fmt.Sprintf("update %s set role=? where planet_id = ? and user_id in (?)", MemberTableName), role, planetId, userIds)
	if err != nil {
		return err
	}

	_, err = p.Exec(append([]any{sql}, args...)...)
	return err
}

func (p *persistence) batchUpdateMembersStatus(planetId int64, userIds []int64, status entity.PlanetMemberStatus) error {
	sql, args, err := sqlx.In(fmt.Sprintf("update %s set status=? where planet_id = ? and user_id in (?)", MemberTableName), status, planetId, userIds)
	if err != nil {
		return err
	}

	_, err = p.Exec(append([]any{sql}, args...)...)
	return err
}

func (p *persistence) getPlanetById(planetId int64) (*entity.Planet, error) {
	var res entity.Planet
	exists, err := p.Table(RootTableName).Where("id = ?", planetId).Get(&res)
	if !exists {
		return nil, nil
	}

	return &res, err
}

func (p *persistence) getMember(planetId, userId int64) (*entity.PlanetMember, error) {
	var res entity.PlanetMember
	exists, err := p.Table(MemberTableName).Where("planet_id = ? and user_id = ?", planetId, userId).Get(&res)
	if !exists {
		return nil, nil
	}

	return &res, err
}

func (p *persistence) createMember(member *entity.PlanetMember) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := p.Table(MemberTableName).Insert(member)
		return err
	})
}

// TODO listPlanetMembersByRole?
func (p *persistence) listAdmins(planetId int64) ([]*entity.PlanetMember, error) {
	var res []*entity.PlanetMember
	return res, p.Table(MemberTableName).Where("planet_id = ? and role >= 1", planetId).Find(&res)
}

func (p *persistence) listPlanetMembersByUserIds(planetId int64, userIds []int64) ([]*entity.PlanetMember, error) {
	var res []*entity.PlanetMember
	return res, p.Table(MemberTableName).In("user_id", userIds).And("planet_id = ?", planetId).Find(&res)
}
