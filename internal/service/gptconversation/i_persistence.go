package gptconversation

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type iPersistence interface {
	create(message *entity.GptChatMessage) error
	updateLikeState(messageId string, likeState entity.LikeState) error
	pageByUserId(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.GptChatMessage], error)
	isEmptyConversation(userId int64) (bool, error)
}
