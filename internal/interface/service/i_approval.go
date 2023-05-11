package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"

type IApproval interface {
	StartApprove(_type entity.ApprovalType, userId int64, reason string, businessId int64) (*entity.Approval, error)
	Audit(id, userId int64, pass bool) error
	GetOne(id int64) (*entity.Approval, bool, error)
	ListByIds(ids []int64) ([]*entity.Approval, error)
}

// BusinessListener 业务监听器，由业务实现
type BusinessListener interface {
	OnPass(approval *entity.Approval) error
	OnReject(approval *entity.Approval) error
}
