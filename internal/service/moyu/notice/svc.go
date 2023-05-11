package notice

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"time"
)

type svc struct {
	gone.Flag
	emitter.Sender `gone:"gone-emitter"`

	P        iPersistence           `gone:"*"`
	Message  service.IMessageRecord `gone:"*"`
	Approval service.IApproval      `gone:"*"`
	User     service.IUser          `gone:"*"`
}

func (s *svc) Create(userId int64, messageType entity.NoticeType, businessId int64, read bool) error {
	message := &entity.Notice{
		UserId:     userId,
		BusinessId: businessId,
		Type:       messageType,
		Read:       read,
		CreatedAt:  time.Now(),
	}

	err := s.P.create(message)
	if err != nil {
		return err
	}

	return s.Send(&entity.NoticeMessageCreatedEvent{Notice: *message})
}

func (s *svc) Get(id int64) (*entity.Notice, bool, error) {
	return s.P.getById(id)
}

func (s *svc) Page(userId int64, query page.Query) (*page.Result[*domain.Notice], error) {
	p, err := s.P.pageByUserId(query, userId)
	if err != nil {
		return nil, err
	}

	var userIds []int64

	linq.From(p.List).Select(func(i interface{}) interface{} {
		return i.(*entity.Notice).UserId
	}).Distinct().ToSlice(&userIds)

	domains, err := s.parseDomain(p.List)
	if err != nil {
		return nil, err
	}

	return &page.Result[*domain.Notice]{List: domains, Total: p.Total}, nil
}

func (s *svc) MarkRead(userId int64, noticeIds []int64) error {
	err := s.P.markRead(userId, noticeIds)
	if err != nil {
		return err
	}

	return s.Send(&entity.NoticeMessageRead{
		UserId: userId,
		Ids:    noticeIds,
	})
}

func (s *svc) CountUnread(userId int64) (int64, error) {
	return s.P.countUnread(userId)
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s *svc) parseBusinessApproval(n []*domain.Notice) error {
	approvalIds := collection.Map(n, func(n *domain.Notice) int64 {
		return n.BusinessId
	})

	approvals, err := s.Approval.ListByIds(approvalIds)
	if err != nil {
		return err
	}

	userIds := collection.Map(approvals, func(a *entity.Approval) int64 {
		return a.UserId
	})

	users, err := s.User.GetUserSimpleInBatch(userIds)
	if err != nil {
		return err
	}

	approvalMap := collection.ToMap(approvals, func(a *entity.Approval) (int64, *entity.Approval) {
		return a.Id, a
	})

	userMap := collection.ToMap(users, func(a *entity.UserSimple) (int64, *entity.UserSimple) {
		return a.Id, a
	})

	for _, d := range n {
		if record, ok := approvalMap[d.BusinessId]; ok {
			if user, ok := userMap[record.UserId]; ok {
				d.User = *user
			}

			d.Business = record
		}
	}

	return nil
}

func (s *svc) parseBusinessMessage(n []*domain.Notice) error {
	messageIds := collection.Map(n, func(n *domain.Notice) int64 {
		return n.BusinessId
	})

	records, err := s.Message.GetRecords(messageIds)
	if err != nil {
		return err
	}

	sourceUserIds := collection.Map(records, func(m *message.Message) int64 {
		return m.UserId
	})

	users, err := s.User.GetUserSimpleInBatch(sourceUserIds)
	if err != nil {
		return err
	}

	userMap := collection.ToMap(users, func(u *entity.UserSimple) (int64, *entity.UserSimple) {
		return u.Id, u
	})

	recordMap := collection.ToMap(records, func(m *message.Message) (int64, *message.Message) {
		return m.Id, m
	})

	for _, d := range n {
		if record, ok := recordMap[d.BusinessId]; ok {
			if user, ok := userMap[record.UserId]; ok {
				d.User = *user
			}

			d.Business = record
		}
	}

	return nil
}

func (s *svc) parseDomain(n []*entity.Notice) ([]*domain.Notice, error) {
	domains := collection.Map(n, func(n *entity.Notice) *domain.Notice {
		return &domain.Notice{Notice: *n}
	})

	var typeGroups []linq.Group
	linq.From(domains).GroupByT(func(n *domain.Notice) entity.NoticeType {
		return n.Type
	}, func(n *domain.Notice) *domain.Notice { return n }).ToSlice(&typeGroups)

	for _, g := range typeGroups {
		switch g.Key.(entity.NoticeType) {
		case entity.NoticeTypeApproval:
			err := s.parseBusinessApproval(collection.Map(g.Group, func(i any) *domain.Notice { return i.(*domain.Notice) }))
			if err != nil {
				return nil, err
			}

		case entity.NoticeTypeMention, entity.NoticeTypeReference:
			err := s.parseBusinessMessage(collection.Map(g.Group, func(i any) *domain.Notice { return i.(*domain.Notice) }))
			if err != nil {
				return nil, err
			}
		}
	}

	return domains, nil

}

func (s *svc) ListUnreadIMMessages(userId, ackMessageId, channelId int64) ([]int64, error) {
	return s.P.listUnreadIMMessage(userId, ackMessageId, channelId)
}

func (s *svc) ListUnreadNotice(noticeType entity.NoticeType, businessId int64) ([]*entity.Notice, error) {
	return s.P.listUnreadNotice(noticeType, businessId)
}
