package identity

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type svc struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s *svc) List() ([]*entity.Identity, error) {
	var res []*entity.Identity
	return res, s.Table("anonymous_identity").Find(&res)
}

func (s *svc) Random() (*entity.Identity, error) {
	i := new(entity.Identity)
	if has, err := s.SQL(`select * from anonymous_identity order by rand() limit 1`).Get(i); err != nil {
		return nil, gin.ToError(err)
	} else if has {
		return i, nil
	} else {
		return nil, nil
	}
}
