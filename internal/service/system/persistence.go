package system

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:gone
func NewResignPersistence() gone.Goner {
	return &persistence{}
}

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *persistence) listByGroupId(group int, sort bool) (list []*entity.Emoticon, err error) {
	query := p.Table("emoticon")
	if sort {
		query = query.Desc("ref_count")
	}
	if group > 0 {
		query = query.Where("`group` = ?", group)
	}

	return list, query.Find(&list)
}

func (p *persistence) updateRefStat() error {
	_, err := p.Exec(`
	update emoticon as a
	inner join (
	select
	a.content->>'$.emoticonId' as emoticonId,
	count(1) as count
	from message_record as a
	where
	a.created_at > date_add(now(),interval -1 day)
	and a.content->>'$.type' = 'emoticon'
	group by emoticonId
	) as b on b.emoticonId = a.id
	set
	a.ref_count = a.ref_count + b.count,
	a.ref_stat_time = now()
`)

	return err
}
