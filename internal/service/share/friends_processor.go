package share

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type FriendsLoginProcessor struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	retain service.IRetainMessage `gone:"*"`
	points service.IPointStrategy `gone:"*"`
	user   service.IUser          `gone:"*"`

	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &FriendsLoginProcessor{}
}

func (e *FriendsLoginProcessor) Consume(on emitter.OnEvent) {
	on(e.handleFirstLogin)
}

func (e *FriendsLoginProcessor) handleFirstLogin(evt *entity.FirstLoginEvent) error {
	e.Logger.Infof("handleFirstLogin: userId=%d", evt.User.Id)

	if evt.InvitedBy != nil {
		invitedBy, err := e.user.GetUserById(*evt.User.InvitedBy)
		if err != nil {
			return err
		}

		if invitedBy == nil {
			return nil
		}

		return e.handleInvitedUserRegister(evt.User, *invitedBy)
	}

	return e.handleNatureUserRegister(evt.User)
}

func (e *FriendsLoginProcessor) processShareUser(invitedBy entity.User, friend entity.User) error {
	points, err := e.points.ApplyPoints(invitedBy.Id, entity.StrategyArgInvite{})
	if err != nil {
		return err
	}

	return e.retain.SendRetainMessage(&entity.RetainMessage{
		Type:   entity.RetainMessageTypeFriendsFirstLogin,
		UserId: invitedBy.Id,
		Payload: entity.RetainMessageFriendsFirstLogin{
			Points: points,
			Friends: []*entity.UserSimple{
				{
					Id:       friend.Id,
					Avatar:   friend.Avatar,
					Nickname: friend.Nickname,
				},
			},
		},
	})
}

func (e *FriendsLoginProcessor) handleInvitedUserRegister(user, invitedBy entity.User) error {
	e.Logger.Infof("handleInvitedUserRegister: userId=%d", user.Id)

	return e.Transaction(func(session xorm.Interface) error {
		if err := e.processShareUser(invitedBy, user); err != nil {
			return err
		}

		if _, err := e.points.ApplyPoints(user.Id, entity.StrategyArgBeInvited{}); err != nil {
			return err
		}

		return nil
	})
}

func (e *FriendsLoginProcessor) handleNatureUserRegister(user entity.User) error {
	e.Logger.Infof("handleNatureUserRegister: userId=%d", user.Id)

	_, err := e.points.ApplyPoints(user.Id, entity.StrategyArgNewRegister{})
	return err
}
