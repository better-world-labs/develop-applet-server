package hotissue

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type svc struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
	Rewriter    *ContentRewriter `gone:"*"`
}

//go:gone
func NewService() gone.Goner {
	return &svc{}
}

func (s *svc) List() ([]*entity.HotIssue, error) {
	var res []*entity.HotIssue
	return res, s.Table("hot_issue").Find(&res)
}

func (s *svc) ListIssues() ([]*entity.HotIssue, error) {
	issues, err := s.List()
	if err != nil {
		return nil, err
	}

	s.Rewriter.RewriteList(issues)
	return issues, nil
}
