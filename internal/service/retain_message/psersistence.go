package retain_message

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type RetainMessageOffset struct {
	UserId   int64
	OffsetId int64
}

type persistence struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

func (p persistence) listRetainMessages(userId, offsetId int64) ([]*entity.RetainMessage, error) {
	var res []*entity.RetainMessage
	return res, p.Where("user_id = ? and id > ?", userId, offsetId).Desc("id").Find(&res)
}

func (p persistence) getReadOffset(userId int64) (int64, error) {
	var res RetainMessageOffset
	_, err := p.Where("user_id = ?", userId).Get(&res)
	return res.OffsetId, err
}

func (p persistence) markOffset(userId, offsetId int64) error {
	_, err := p.Exec("insert retain_message_offset (user_id, offset_id) values(?, ?)"+
		"on duplicate key update offset_id = ?", userId, offsetId, offsetId)
	return err
}

func (p persistence) create(message *entity.RetainMessage) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.Insert(message)
		return err
	})
}
