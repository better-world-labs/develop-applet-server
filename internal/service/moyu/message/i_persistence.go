package message

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type (
	IPersistence interface {
		createOrUpdateLike(like *entity.MessageLike) error

		create(record *message.Message) error

		checkExists(sendId string) (bool, error)

		listHistoryMessageBefore(channelId, fromId int64, size int) ([]*message.Message, error)

		listHistoryMessageAfter(channelId, fromId int64, size int) ([]*message.Message, error)

		GetLastMessageByChannelId(channelId int64) (*message.Message, error)

		countMessageLikes(messageIds []int64) ([]*entity.MessageLikeCount, error)

		listRangedMessageIdsByUserId(messageIds []int64, userId int64) ([]int64, error)

		getById(id int64) (*message.Message, bool, error)

		listByIds(ids []int64) ([]*message.Message, error)

		listByIdsMap(ids []int64) (map[int64]*message.Message, error)

		GetRecordsSummary(fromId, channelId int64) (int64, error)

		GetMessageCntByUserId(userId int64) (int64, error)

		insertReplyForMessage(msgId int64, emoticonId int64, userId int64)
	}
)
