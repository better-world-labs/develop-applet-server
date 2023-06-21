package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IJSSDK interface {
	GetSignature(url string) (entity.JSSDKSignature, error)
}
