package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type ISystem interface {
	GenOssToken(userId int64, fileType string) (*entity.OssToken, error)

	//GetEmoticonList 获取表情列表
	GetEmoticonList(group int, sort bool) ([]*entity.Emoticon, error)
}
