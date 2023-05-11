package channel

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/moyu/approval"
)

type listener struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	svc    service.IChannel `gone:"*"`
	Sender emitter.Sender   `gone:"gone-emitter"`
}

//go:gone
func NewListener() gone.Goner {
	l := &listener{}
	approval.RegisterBusinessListener(entity.ApprovalTypeChannelJoin, l)
	return l
}

func (l *listener) OnPass(approval *entity.Approval) error {
	l.Infof("onApprovalPass: businessId=%d", approval.BusinessId)
	err := l.svc.UpdateMemberState(approval.UserId, approval.BusinessId, entity.ChannelMemberStateJoined)
	if err != nil {
		return err
	}

	//TODO 封装在 channelMember 里
	return l.Sender.Send(&entity.ChannelMemberAdded{ChannelId: approval.BusinessId, UserId: approval.UserId})
}

func (l *listener) OnReject(approval *entity.Approval) error {
	l.Infof("onApprovalReject: businessId=%d", approval.BusinessId)
	return l.svc.UpdateMemberState(approval.UserId, approval.BusinessId, entity.ChannelMemberStateNotJoin)
}
