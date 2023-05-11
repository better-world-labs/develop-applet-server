package stat

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewMessageStat() gone.Goner {
	return &messageStat{}
}

type messageStat struct {
	gone.Flag
	p        iPersistence           `gone:"*"`
	iMessage service.IMessageRecord `gone:"*"`
}

func (s *messageStat) HotMessageTop(channelId int64, top int) (out []*entity.HotMessage, err error) {
	if top <= 0 {
		top = 10
	}
	if top > 50 {
		top = 50
	}

	list, err := s.p.listTopReplyMsg(top, channelId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	msgIds := collection.Map(list, func(msg *hotMsg) int64 {
		return msg.MsgId
	})

	recordsMap, err := s.iMessage.GetRecordsMap(msgIds)
	if err != nil {
		return nil, gin.ToError(err)
	}

	out = make([]*entity.HotMessage, 0, len(list))
	for _, it := range list {
		record := recordsMap[it.MsgId]
		if record != nil {
			if err != nil {
				return nil, gin.ToError(err)
			}
			out = append(out, &entity.HotMessage{
				MsgId:      record.Id,
				UserId:     record.UserId,
				ReplyCount: it.ReplyCount,
				Content:    record.Content,
			})
		}
	}
	return
}
