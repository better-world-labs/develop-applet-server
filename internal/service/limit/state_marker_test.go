package limit

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
)

func TestStateMarker(t *testing.T) {
	gone.Test(func(p *marker) {
		sent, err := p.IsAppCreateNotifySent(10069)
		if err != nil {
			t.Fatalf("%v\n", err)
			return
		}

		err = p.MarkAppCreatePointsLimitNotifySent(10069)
		if err != nil {
			t.Fatalf("%v\n", err)
			return
		}

		sent, err = p.IsAppCreateNotifySent(10069)
		if err != nil {
			t.Fatalf("%v\n", err)
			return
		}

		t.Logf("sent: %v\n", sent)
	}, func(cemetery gone.Cemetery) error {
		_ = test.RedisPriest(cemetery)
		cemetery.Bury(NewMarker())
		return nil
	})
}
