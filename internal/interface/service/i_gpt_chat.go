package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IGPTChat interface {
	SendMessage(userId int64, content string) (ChannelStreamTrunkReader[entity.GptChatMessage], error)
	ListMessages(userId int64) ([]*entity.GptChatMessage, error)
}
