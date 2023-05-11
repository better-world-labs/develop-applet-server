package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IAnonymousIdentity interface {
	List() ([]*entity.Identity, error)

	//Random 随机获取一个身份
	Random() (*entity.Identity, error)
}
