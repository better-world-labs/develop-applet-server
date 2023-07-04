package gptconversation

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
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

func (p *persistence) create(message *entity.GptChatMessage) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.Insert(message)
		return err
	})
}

func (p *persistence) pageByUserId(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.GptChatMessage], error) {
	var res []*entity.GptChatMessage

	session := p.Where("user_id = ?", userId)
	if query.CursorIndicator() > 0 {
		session.Where("id < ?", query.CursorIndicator())
	}

	if err := session.Desc("id").Limit(query.Size(), 0).Find(&res); err != nil {
		return nil, err
	}

	return page.NewStreamResult(res), nil
}

func (p *persistence) isEmptyConversation(userId int64) (bool, error) {
	has, err := p.Where("user_id = ?", userId).Exist()
	return !has, err
}
