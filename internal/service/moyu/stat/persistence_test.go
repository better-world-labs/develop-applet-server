package stat

import (
	"github.com/gone-io/gone"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
)

func Test_persistence_topReplyMsg(t *testing.T) {
	gone.Test(func(p *persistence) {
		list, err := p.listTopReplyMsg(10, 42)
		assert.Nil(t, err)

		assert.True(t, len(list) > 0)

	}, func(cemetery gone.Cemetery) error {
		_ = test.MysqlPriest(cemetery)
		cemetery.Bury(NewPersistence())
		return nil
	})
}
