package entity

import (
	"github.com/gone-io/gone/goner/gin"
	"time"
)

type (
	ApprovalType  string
	ApprovalState uint
)

const (
	ApprovalTypeChannelJoin ApprovalType  = "channel-join"
	ApprovalStateTodo       ApprovalState = 0
	ApprovalStatePass       ApprovalState = 1
	ApprovalStateReject     ApprovalState = 2
)

var (
	ErrorAlreadyApproved = gin.NewBusinessError("已经审批过了")
)

type Approval struct {
	Id           int64         `json:"id"`
	ApprovalType ApprovalType  `json:"approvalType"`
	Reason       string        `json:"reason"`
	BusinessId   int64         `json:"businessId"`
	Business     any           `xorm:"-" json:"business"`
	UserId       int64         `json:"userId"`
	CreatedAt    time.Time     `json:"createdAt"`
	AuditBy      *int64        `json:"audit_by"`
	AuditAt      *time.Time    `json:"auditAt"`
	State        ApprovalState `json:"state"`
}

func NewApproval(approvalType ApprovalType, reason string, businessId int64, userId int64, cratedAt time.Time) *Approval {
	return &Approval{
		ApprovalType: approvalType,
		Reason:       reason,
		BusinessId:   businessId,
		UserId:       userId,
		CreatedAt:    cratedAt,
	}
}

func (a *Approval) Audit(userId int64, pass bool) error {
	if a.State != ApprovalStateTodo {
		return ErrorAlreadyApproved
	}

	t := time.Now()
	a.AuditAt = &t
	a.AuditBy = &userId

	if pass {
		a.State = ApprovalStatePass
	} else {
		a.State = ApprovalStateReject
	}

	return nil
}
