package miniapp

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/require"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
)

func Priests(cemetery gone.Cemetery) error {
	_ = goner.BasePriest(cemetery)
	_ = test.MysqlPriest(cemetery)
	cemetery.Bury(NewPMiniApp())
	return nil
}

func TestPersistenceMiniAppGetLastNOutputs(t *testing.T) {
	gone.Test(func(p *pMiniApp) {
		outputs, err := p.getLastNOutputByAppIds([]string{"1111", "2222"}, 2)
		require.Nil(t, err)

		t.Logf("outputs: %v\n", outputs)
	}, Priests)
}

func TestPersistence(t *testing.T) {
	gone.Test(func(p *pMiniApp) {
		count, err := p.countAppByUserId(10068)
		require.Nil(t, err)

		t.Logf("count: %v\n", count)

		count, err = p.countOutputsAppIdByUserId(10068)
		require.Nil(t, err)

		t.Logf("count: %v\n", count)
	}, Priests)
}
