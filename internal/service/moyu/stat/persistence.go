package stat

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
)

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *persistence) listTopReplyMsg(top int, channelId int64) (list []*hotMsg, err error) {
	return list, p.SQL(`
		select
			a.i_msg_ref as msg_id,
			count(1) as reply_count
		from message_record as a
		where
			a.i_msg_ref > 0
			and a.channel_id = ?
		group by msg_id
		order by reply_count desc
		limit ?
    `, channelId, top).Find(&list)
}
