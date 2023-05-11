package approval

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"time"
)

var listeners = map[entity.ApprovalType]service.BusinessListener{}

func RegisterBusinessListener(approvalType entity.ApprovalType, listener service.BusinessListener) {
	if _, exists := listeners[approvalType]; exists {
		panic(fmt.Sprintf("BusinessListener for %s already exists", approvalType))
	}

	listeners[approvalType] = listener
}

type svc struct {
	gone.Flag
	p              iPersistence `gone:"*"`
	emitter.Sender `gone:"gone-emitter"`
	xorm.Engine    `gone:"gone-xorm"`
	Notice         service.INotice  `gone:"*"`
	Channels       service.IChannel `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s *svc) StartApprove(_type entity.ApprovalType, userId int64, reason string, businessId int64) (*entity.Approval, error) {
	approval := entity.NewApproval(_type, reason, businessId, userId, time.Now())

	return approval, s.Transaction(func(session xorm.Interface) error {
		approval := approval
		err := s.p.create(approval)
		if err != nil {
			return err
		}

		return s.Send(&entity.ApprovalStarted{Approval: *approval})
	})
}

func (s *svc) ListByIds(ids []int64) ([]*entity.Approval, error) {
	approvals, err := s.p.listByIds(ids)
	if err != nil {
		return nil, err
	}

	return approvals, s.processBusiness(approvals)
}

func (s *svc) GetOne(id int64) (*entity.Approval, bool, error) {
	return s.p.getById(id)
}

func (s *svc) Audit(id, userId int64, pass bool) error {
	return s.Transaction(func(session xorm.Interface) error {
		approval, exists, err := s.p.getById(id)
		if err != nil {
			return err
		}

		if !exists {
			return gin.NewParameterError("approval not found")
		}

		err = approval.Audit(userId, pass)
		if err != nil {
			return err
		}

		err = s.p.update(approval)
		if err != nil {
			return err
		}

		if approval.State == entity.ApprovalStatePass {
			err = listeners[approval.ApprovalType].OnPass(approval)
			if err != nil {
				return err
			}
		} else {
			err = listeners[approval.ApprovalType].OnReject(approval)
			if err != nil {
				return err
			}
		}

		return s.Send(&entity.ApprovalAudited{Id: id})
	})
}

func (s *svc) processBusiness(approvals []*entity.Approval) error {
	var typeGroups []linq.Group
	linq.From(approvals).GroupByT(func(n *entity.Approval) entity.ApprovalType {
		return n.ApprovalType
	}, func(n *entity.Approval) *entity.Approval { return n }).ToSlice(&typeGroups)

	for _, g := range typeGroups {
		switch g.Key.(entity.ApprovalType) {
		case entity.ApprovalTypeChannelJoin:
			err := s.processBusinessChannel(collection.Map(g.Group, func(i any) *entity.Approval { return i.(*entity.Approval) }))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *svc) processBusinessChannel(approvals []*entity.Approval) error {
	businessIds := collection.Map(approvals, func(approval *entity.Approval) int64 {
		return approval.BusinessId
	})

	channels, err := s.Channels.ListChannelsByIds(businessIds)
	if err != nil {
		return err
	}

	channelMap := collection.ToMap(channels, func(channel *entity.Channel) (int64, *entity.Channel) {
		return channel.Id, channel
	})

	for _, a := range approvals {
		a.Business = channelMap[a.BusinessId]
	}

	return nil
}
