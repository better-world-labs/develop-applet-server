package push

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPushMessage(t *testing.T) {
	gone.Test(func(p *PushSvc) {
		err := p.PushMessage(10067, map[string]any{
			"type": "share-hint-use-app",
			"payload": map[string]any{
				"usedApps":   1,
				"costPoints": 12,
			},
		})
		assert.Nil(t, err)

	}, func(cemetery gone.Cemetery) error {
		_ = goner.BasePriest(cemetery)
		cemetery.Bury(NewSvc())
		return nil
	})
}
