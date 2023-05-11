package channel

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"strconv"
	"strings"
)

const (
	GroupTableName        = "channel_group"
	RootTableName         = "channel"
	SubscriptionTableName = "channel_member"
)

//go:gone
func NewChanelPersistence() gone.Goner {
	return &channelPersistence{}
}

//go:gone
func NewChanelGroupPersistence() gone.Goner {
	return &channelGroupPersistence{}
}

//go:gone
func NewChanelMemberPersistence() gone.Goner {
	return &channelMemberPersistence{}
}

type channelPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

type channelGroupPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

type channelMemberPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *channelGroupPersistence) delete(id int64) error {
	channel := new(entity.ChannelGroup)
	_, err := p.Table(GroupTableName).ID(id).Delete(channel)
	return err
}

func (p *channelGroupPersistence) listByPlanetId(planetId int64) ([]*entity.ChannelGroup, error) {
	res := make([]*entity.ChannelGroup, 0)
	err := p.Table(GroupTableName).Select("id, name, icon").Where("planet_id = ?", planetId).Asc("sort").Find(&res)
	return res, err
}

func (p *channelGroupPersistence) create(group *entity.ChannelGroup) error {
	_, err := p.Table(GroupTableName).Insert(group)
	return err
}

func (p *channelGroupPersistence) updateSort(planetId, channelId, sort int64) error {
	channel := entity.ChannelGroup{
		Id:       channelId,
		PlanetId: planetId,
		Sort:     sort,
	}
	_, err := p.Table(GroupTableName).ID(channelId).Cols("sort").Update(channel)
	return err
}
func (p *channelGroupPersistence) get(groupId int64) (*entity.ChannelGroup, error) {
	var res entity.ChannelGroup
	exist, err := p.Table(GroupTableName).Where("id = ?", groupId).Get(&res)
	if err != nil {
		return nil, err
	}

	if exist {
		return &res, err
	}

	return nil, err
}

func (p *channelGroupPersistence) updateName(channelGroupId int64, name string) error {
	group := entity.ChannelGroup{Name: name}
	_, err := p.Table(GroupTableName).Where("id = ?", channelGroupId).Cols("name").Update(group)
	return err
}

func (p *channelPersistence) updateSort(groupId, channelId, sort int64) error {
	channel := entity.Channel{
		Id:      channelId,
		GroupId: groupId,
		Sort:    sort,
	}
	_, err := p.Table(RootTableName).ID(channelId).Cols("sort").Update(channel)
	return err
}

func (p *channelPersistence) listByPlanetId(planetId int64) ([]*entity.Channel, error) {
	res := make([]*entity.Channel, 0)
	err := p.Table(RootTableName).Where("planet_id = ? and state = 1", planetId).OrderBy("sort asc, id asc").Find(&res)
	return res, err
}

func (p *channelPersistence) create(channel *entity.Channel) error {
	_, err := p.Table(RootTableName).Insert(channel)
	return err
}

func (p *channelPersistence) delete(channelId int64) error {
	aChannel := new(entity.Channel)
	_, err := p.Table(RootTableName).ID(channelId).Delete(aChannel)
	return err
}

func (p *channelPersistence) get(channelId int64) (*entity.Channel, bool, error) {
	var res entity.Channel
	exist, err := p.Table(RootTableName).Where("id = ?", channelId).Get(&res)
	return &res, exist, err
}

func (p *channelPersistence) listByIds(ids []int64) ([]*entity.Channel, error) {
	var res []*entity.Channel
	return res, p.Table(RootTableName).In("id", ids).Find(&res)
}

func (p *channelPersistence) updateNotice(channelId int64, notice string) error {
	group := entity.Channel{Notice: notice}
	_, err := p.Table(RootTableName).Where("id = ?", channelId).Cols("notice").Update(group)
	return err
}

func (p *channelPersistence) updateName(channelId int64, name string) error {
	group := entity.Channel{Name: name}
	_, err := p.Table(RootTableName).Where("id = ?", channelId).Cols("name").Update(group)
	return err
}

func (p *channelMemberPersistence) list(channelId int64) ([]*entity.ChannelMember, error) {
	var res []*entity.ChannelMember
	err := p.Table(SubscriptionTableName).Where("channel_id = ? and state = ?", channelId, entity.ChannelMemberStateJoined).Find(&res)
	return res, err
}

func (p *channelMemberPersistence) listIdsByChannelId(channelId int64) ([]int64, error) {
	res := make([]int64, 0)
	err := p.Table(SubscriptionTableName).Select("user_id").Where("channel_id = ? and state = ?", channelId, entity.ChannelMemberStateJoined).Find(&res)
	return res, err
}

func (p *channelMemberPersistence) listMembersOffset(channelId int64, userIds []int64) (map[int64]int64, error) {
	var arr []*entity.MessageOffset

	err := p.Table(SubscriptionTableName).Cols("user_id", "channel_id", "last_read_message_id").
		Where("channel_id = ? ", channelId).
		In("user_id", userIds).Find(&arr)

	return collection.ToMap(arr, func(i *entity.MessageOffset) (int64, int64) {
		return i.UserId, i.LastReadMessageId
	}), err
}

func (p *channelMemberPersistence) listChannelIdsByUserId(userId int64) ([]int64, error) {
	var ids []int64
	return ids, p.Table(SubscriptionTableName).
		Cols("channel_id").
		Where("user_id = ?", userId).Find(&ids)
}

func (p *channelMemberPersistence) isExists(channelId, userId int64) (bool, error) {
	return p.Table(SubscriptionTableName).
		Where("channel_id = ? and user_id = ? and state = ?", channelId, userId, entity.ChannelMemberStateJoined).Exist()
}

func (p *channelMemberPersistence) create(member *entity.ChannelMember) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := p.Exec(`
	insert channel_member (channel_id, user_id, last_read_message_id, created_at, state, apply_id)
	select ?, ?, id, ?, ?, ?
	from ( select id from message_record where channel_id = ? union select 0 order by id desc
	limit 1) a
	on duplicate key update apply_id = ?, state = ?, created_at = ?`,
			member.ChannelId,
			member.UserId,
			member.CreatedAt,
			member.State,
			member.ApplyId,
			member.ChannelId,
			member.ApplyId,
			member.State,
			member.CreatedAt,
		)
		return err
	})
}

func (p *channelMemberPersistence) delete(id int64) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := p.Table(SubscriptionTableName).Where("id = ?", id).Delete()
		return err
	})
}

func (p *channelMemberPersistence) get(channelId, userId int64) (*entity.ChannelMember, error) {
	var res entity.ChannelMember
	exists, err := p.Table(SubscriptionTableName).
		Where("channel_id = ? and user_id = ?", channelId, userId).Get(&res)
	if !exists {
		return nil, nil
	}

	return &res, err
}

func (p *channelMemberPersistence) updateState(id int64, state entity.ChannelMemberState) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := p.Exec("update channel_member set state = ? where id = ?", state, id)
		return err
	})
}

func (p *channelMemberPersistence) UpdateLastReadMsgId(userId, channelId, msgId int64) error {
	member := entity.ChannelMember{
		LastReadMessageId: msgId,
	}
	// field last_read_message_id must be increase
	_, err := p.Table(SubscriptionTableName).Where("user_id  = ? and channel_id = ? and last_read_message_id < ? ", userId, channelId, msgId).Cols("last_read_message_id").Update(member)

	return err
}

func (p *channelMemberPersistence) BatchUpdateLastReadMsgId(userIds []int64, channelId, msgId int64) error {
	var strIds []string
	linq.From(userIds).SelectT(func(id int64) string {
		return strconv.FormatInt(id, 10)
	}).ToSlice(&strIds)

	// field last_read_message_id must be increase
	sql := fmt.Sprintf("update %s set last_read_message_id = ? where channel_id = ? and user_id in (%s)", SubscriptionTableName, strings.Join(strIds, ","))

	_, err := p.Exec(sql, msgId, channelId)
	return err
}

func (p *channelMemberPersistence) GetChannelIdsByUserId(userId int64) ([]int64, error) {
	var channelIds []int64

	err := p.Table(SubscriptionTableName).Select("channel_id").Where("user_id  = ? ", userId).Find(&channelIds)

	return channelIds, err
}
