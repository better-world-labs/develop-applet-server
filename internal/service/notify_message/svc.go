package notify_message

import (
	"fmt"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"time"
)

const TableName = "notify_message"

type svc struct {
	gone.Flag
	emitter.Sender `gone:"gone-emitter"`
	xorm.Engine    `gone:"gone-xorm"`
}

//go:gone
func NewSvs() gone.Goner {
	return &svc{}
}

func (s svc) SendNotify(_type entity.NotifyMessageType, userId int64, Title, content string, operationId string) error {
	return s.Send(&entity.CreateNotifyMessageEvent{
		NotifyMessageInfo: entity.NotifyMessageInfo{
			Type:        _type,
			UserId:      userId,
			Title:       Title,
			Content:     content,
			OperationId: operationId,
		},
	})
}

func (s svc) checkNotifyMessageExistsByOperationId(operationId string) (bool, error) {
	return s.Table(entity.NotifyMessage{}).Where("operation_id = ?", operationId).Exist()
}

func (s svc) sendNotifySync(message entity.NotifyMessageInfo) error {
	return s.Transaction(func(session xorm.Interface) error {
		message := &entity.NotifyMessage{
			NotifyMessageInfo: message,
			CreatedAt:         time.Now(),
		}

		err := s.createNotifyMessage(message)
		if err != nil {
			return err
		}

		return s.Send(&entity.NotifyMessageCreatedEvent{
			NotifyMessage: *message,
		})
	})
}

func (s svc) createNotifyMessage(message *entity.NotifyMessage) error {
	return s.Transaction(func(session xorm.Interface) error {
		err := message.Validate()
		if err != nil {
			return err
		}

		_, err = session.Insert(message)
		return err
	})
}

func (s svc) PageNotifyMessages(userId int64, query page.StreamQuery, filter entity.NotifyMessageListFilter) (*page.StreamResult[*entity.NotifyMessage], error) {
	var slice []*entity.NotifyMessage

	session := s.Where("user_id = ?", userId)

	if query.CursorIndicator() > 0 {
		session = session.Where("id < ?", query.CursorIndicator())
	}

	if filter.IsRead != nil {
		session = session.Where("is_read = ?", *filter.IsRead).MustCols("is_read")
	}

	err := session.Desc("id").Limit(query.Size(), 0).Find(&slice)
	if err != nil {
		return nil, err
	}

	return page.NewStreamResult[*entity.NotifyMessage](slice), nil
}

func (s svc) CountUnread(userId int64) (int64, error) {
	return s.Table(entity.NotifyMessage{}).Where("user_id = ? and is_read = 0", userId).Count()
}

func (s svc) MarkRead(userId, id int64) error {
	_, err := s.Exec(fmt.Sprintf("update %s set is_read = 1 where id = ? and user_id = ?", TableName), id, userId)
	return err
}

func (s svc) MarkReadAll(userId int64) error {
	_, err := s.Exec(fmt.Sprintf("update %s set is_read = 1 where user_id = ?", TableName), userId)
	return err
}
