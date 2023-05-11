package notify_message

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/jmoiron/sqlx"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

func (p *persistence) pageByUserId(query page.Query, userId int64) (*page.Result[*entity.Notice], error) {
	var result page.Result[*entity.Notice]

	count, err := p.Table(TableName).Where("user_id = ?", userId).Asc("is_read").Desc("created_at").Limit(query.LimitOffset(), query.LimitStart()).FindAndCount(&result.List)
	if err != nil {
		return &result, err
	}

	result.Total = count
	return &result, err
}

func (p *persistence) create(message *entity.Notice) error {
	_, err := p.Table(TableName).Insert(message)
	return err
}

func (p *persistence) markRead(userId int64, messageIds []int64) error {
	sql, args, err := sqlx.In(fmt.Sprintf("update %s set `is_read` = 1 where user_id = ? and `id` in (?) ", TableName), userId, messageIds)
	if err != nil {
		return err
	}

	args = append([]any{sql}, args...)
	_, err = p.Exec(args...)
	return err
}

func (p *persistence) countUnread(userId int64) (int64, error) {
	return p.Table(TableName).Where("user_id = ? and `is_read` = 0", userId).Count()
}

func (p *persistence) listUnreadNotice(noticeType entity.NoticeType, businessId int64) ([]*entity.Notice, error) {
	var slice []*entity.Notice
	return slice, p.Table(TableName).Where("business_id = ? and type = ? and is_read = 0", businessId, noticeType).Find(&slice)
}

func (p *persistence) listUnreadIMMessage(userId, ackMessageId, channelId int64) ([]int64, error) {
	var slice []*entity.Notice

	err := p.SQL(`select n.id id from(
	select * from notice where
	user_id = ? and
	is_read = 0 and
	type in ('reference', 'mention') ) n left join
	message_record m on
	n.business_id = m.id
	where
	m.channel_id = ? and m.id <= ?`, userId, channelId, ackMessageId).Find(&slice)

	return collection.Map(slice, func(n *entity.Notice) int64 {
		return n.Id
	}), err
}

func (p *persistence) getById(id int64) (*entity.Notice, bool, error) {
	var n *entity.Notice
	exists, err := p.Table(TableName).Where("id = ?", id).Get(&n)
	return n, exists, err
}
