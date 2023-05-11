package trigger

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type ApprovalHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender   emitter.Sender    `gone:"gone-emitter"`
	Approval service.IApproval `gone:"*"`
}

//go:gone
func NewApprovalHandler() gone.Goner {
	return &ApprovalHandler{}
}

func (e *ApprovalHandler) Consume(on emitter.OnEvent) {
	on(e.handleApproval)
}

func (e *ApprovalHandler) handleApproval(approval *entity.ApprovalAudited) error {
	e.Infof("handleApproval: id=%d", approval.Id)

	one, exists, err := e.Approval.GetOne(approval.Id)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	return e.Sender.Send(&wsevent.TriggerEvent{
		Scope:  wsevent.TriggerEventScopeByUser,
		Type:   TriggerApprovalAudited,
		UserId: one.UserId,
		Params: []interface{}{one},
	})
}
