package entity

type (
	// ApprovalStarted 一个审批发起时触发
	ApprovalStarted struct {
		Approval
	}

	// ApprovalAudited 一个审批单被审批时触发
	ApprovalAudited struct {
		Id int64 `json:"id"`
	}
)
