package resign

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
)

//go:gone
func NewResignService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Flag
	P iPersistence `gone:"*"`
}

func (s *svc) ListTemplates() (*domain.ResignList, error) {
	templates, err := s.P.listResignTemplates()
	if err != nil {
		return nil, err
	}

	return &domain.ResignList{ResignList: templates}, nil
}
