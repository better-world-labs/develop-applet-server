package message

import (
	"encoding/json"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	event "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"time"
)

type svc struct {
	gone.Flag
	Sender      emitter.Sender `gone:"gone-emitter"`
	xorm.Engine `gone:"gone-xorm"`

	Dao IPersistence `gone:"*"`

	ChannelService service.IChannel        `gone:"*"`
	Context        service.IContextManager `gone:"*"`

	logrus.Logger `gone:"gone-logger"`
}

//go:gone
func NewService() gone.Goner {
	return &svc{}
}

// GetContext 读取上下文， 新版 AI Robot 上线 后改为 message.Message,并接到通用接口
func (s *svc) GetContext(messageId, channelId int64) ([]*wsevent.Msg, error) {
	var slice []*wsevent.Msg
	var done bool
	messages, err := s.Context.GetContext(channelId)
	if err != nil {
		return nil, err
	}

	linq.From(messages).OrderByT(func(message *message.Message) int64 {
		return message.Id
	}).SelectT(func(msg *message.Message) *wsevent.Msg {
		content, err := json.Marshal(msg.Content)
		if err != nil {
			return nil
		}

		// TODO 统一换成 message.Message
		return &wsevent.Msg{
			Id:        msg.Id,
			CreatedAt: msg.CreatedAt,
			UserId:    msg.UserId,
			ChannelId: msg.ChannelId,
			Content:   content,
		}
	}).WhereT(func(msg *wsevent.Msg) bool { return msg != nil }).ForEachT(func(msg *wsevent.Msg) {
		if msg.Id == messageId {
			done = true
		}

		if !done {
			slice = append(slice, msg)
		}
	})

	return slice, nil
}

func (s *svc) SendMessage(clientId string, userId, channelId int64, content []byte) error {
	now := time.Now()
	return s.Sender.Send(&wsevent.MsgSendEvent{
		ClientId: clientId,
		Id:       uuid.NewString(),
		UserId:   userId,
		SeqId:    0,
		Msg: wsevent.Msg{
			CreatedAt: now,
			UserId:    userId,
			ChannelId: channelId,
			Content:   content,
		},
		CreatedAt: now,
	})
}

func (s *svc) SaveMessage(message *message.Message) error {
	message.CreatedAt = time.Now()
	return s.Transaction(func(session xorm.Interface) error {
		err := s.Dao.create(message)
		if err != nil {
			return err
		}
		return s.Context.WriteMessage(message)
	})
}

func (s *svc) ListHistory(channelId, fromId int64, size int, upFlag bool) ([]*message.Message, error) {
	s.Debugf("ListHistory: channelId=%d, fromId=%d, size=%s", channelId, fromId, size)
	var res []*message.Message
	var err error

	_, exist, err := s.ChannelService.GetChannelByChannelId(channelId)
	if err != nil {
		s.Logger.Errorf("Get channel [%d] msg failed, err: %s", channelId, err.Error())
		return res, err
	}

	if !exist {
		s.Logger.Infof("channel [%d] has been deleted, return empty messages. ", channelId)
		return res, err
	}

	if upFlag {
		res, err = s.Dao.listHistoryMessageBefore(channelId, fromId, size)
	} else {
		res, err = s.Dao.listHistoryMessageAfter(channelId, fromId, size)
	}
	s.Debugf("ListHistory: len=%s", len(res))
	return res, err

}

func (s *svc) GetLastMessageByChannelId(channelId int64) (*message.Message, error) {
	return s.Dao.GetLastMessageByChannelId(channelId)
}

func (s *svc) Like(messageId, userId int64, isLike bool) error {
	return s.Sender.Send(&event.MessageLikeEvent{
		MessageLike: event.MessageLike{
			MessageId: messageId,
			UserId:    userId,
			IsLike:    isLike,
		},
	})
}

func (s *svc) CreateLikeMessage(like *event.MessageLike) error {
	return s.Dao.createOrUpdateLike(like)
}

func (s *svc) GetMessageLikes(userId int64, messageIds []int64) ([]*domain.MessageLike, error) {
	slice := make([]*domain.MessageLike, 0, len(messageIds))

	messages, err := s.Dao.countMessageLikes(messageIds)
	if err != nil {
		return nil, err
	}

	messagesMap := collection.ToMap(messages, func(i *event.MessageLikeCount) (int64, *event.MessageLikeCount) {
		return i.Id, i
	})

	likes, err := s.Dao.listRangedMessageIdsByUserId(messageIds, userId)
	if err != nil {
		return nil, err
	}

	set := collection.NewSet[int64]()
	set.AddAll(likes)

	for _, id := range messageIds {
		item := &domain.MessageLike{
			MessageLikeCount: event.MessageLikeCount{
				Id: id,
			},
		}

		if m, ok := messagesMap[id]; ok {
			item.MessageLikeCount = *m
			if set.Contains(id) {
				item.IsLike = true
			}
		}

		slice = append(slice, item)
	}

	return slice, nil
}

func (s *svc) GetRecord(id int64) (*message.Message, bool, error) {
	return s.Dao.getById(id)
}

func (s *svc) GetRecords(ids []int64) ([]*message.Message, error) {
	return s.Dao.listByIds(ids)
}

func (s *svc) GetRecordsMap(ids []int64) (map[int64]*message.Message, error) {
	return s.Dao.listByIdsMap(ids)
}

func (s *svc) GetRecordsNum(fromId, channelId int64) (int64, error) {
	return s.Dao.GetRecordsSummary(fromId, channelId)
}

func (s *svc) MsgReplyRecord(msgId int64, emoticonId int64, userId int64) {
	s.Dao.insertReplyForMessage(msgId, emoticonId, userId)
}

func (s *svc) GetMessageCntByUserId(userId int64) (int64, error) {
	return s.Dao.GetMessageCntByUserId(userId)
}
