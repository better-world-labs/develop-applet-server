package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IContentAudit interface {

	// ScanText 文本内容同步审核，审核通过 CheckPass 返回 true
	ScanText(text string) (*entity.AuditResult, error)

	// ScanImage 图片内容同步审核，审核通过 CheckPass返回 true
	ScanImage(imgUrl string) (*entity.AuditResult, error)
}
