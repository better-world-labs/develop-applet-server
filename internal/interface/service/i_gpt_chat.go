package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type IGPTChat interface {
	SendMessage(userId int64, content string) (*ChannelStreamTrunkReader[entity.GptChatMessage], error)
	ListMessages(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.GptChatMessage], error)
}
