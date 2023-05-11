package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IMessageRecord interface {
	// GetContext 读取缓存的聊天会话窗口，窗口大小受系统配置限制
	GetContext(messageId, channelId int64) ([]*wsevent.Msg, error)

	SendMessage(clientId string, userId, channelId int64, content []byte) error

	SaveMessage(message *message.Message) error

	CreateLikeMessage(like *entity.MessageLike) error

	Like(messageId, userId int64, isLike bool) error

	GetMessageLikes(userId int64, messageIds []int64) ([]*domain.MessageLike, error)

	ListHistory(channelId, fromId int64, size int, upFlag bool) ([]*message.Message, error)

	GetLastMessageByChannelId(channelId int64) (*message.Message, error)

	GetRecord(id int64) (*message.Message, bool, error)

	GetRecords(ids []int64) ([]*message.Message, error)

	GetRecordsMap(ids []int64) (map[int64]*message.Message, error)

	GetRecordsNum(fromId, channelId int64) (int64, error)

	GetMessageCntByUserId(userId int64) (int64, error)

	//MsgReplyRecord 消息被回复了 "表情回复" 消息，更新该消息的reply字段
	MsgReplyRecord(msgId int64, emoticonId int64, userId int64)
}
