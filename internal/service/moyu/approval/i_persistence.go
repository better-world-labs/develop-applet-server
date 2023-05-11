package approval

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type iPersistence interface {
	listByIds(ids []int64) ([]*entity.Approval, error)

	getById(id int64) (*entity.Approval, bool, error)

	create(approval *entity.Approval) error

	update(approval *entity.Approval) error

	delete(id int64) error
}
