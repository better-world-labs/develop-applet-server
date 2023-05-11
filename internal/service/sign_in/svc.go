package sign_in

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"time"
)

type svc struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`

	points service.IPointStrategy `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) GetSignInStatus(userId int64) (bool, error) {
	_, has, err := s.getSign(userId, time.Now())
	return has, err
}

func (s svc) SignIn(userId int64) error {
	status, err := s.GetSignInStatus(userId)
	if err != nil {
		return err
	}

	if status {
		return businesserrors.ErrorAlreadySignIn
	}

	if err := s.createSign(SignInDaily{
		UserId: userId,
		Date:   time.Now(),
	}); err != nil {
		return err
	}

	_, err = s.points.ApplyPoints(userId, entity.StrategyArgSignInDaily{})
	return err

}

func (s svc) createSign(sign SignInDaily) error {
	return s.Transaction(func(session xorm.Interface) error {
		_, err := session.Insert(sign)
		return err
	})
}

func (s svc) getSign(userId int64, date time.Time) (SignInDaily, bool, error) {
	var sign SignInDaily
	exists, err := s.Where("user_id = ? and date = ?", userId, date.Format("2006-01-02")).Get(&sign)
	return sign, exists, err
}
