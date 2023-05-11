package resign

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:generate sh -c "mockgen -package=$GOPACKAGE -self_package=moyu-server/internal/service/$GOPACKAGE  -source=$GOFILE|gone mock -o persistence_mock_test.go"
type iPersistence interface {
	listResignTemplates() ([]*entity.ResignTemplate, error)
}
