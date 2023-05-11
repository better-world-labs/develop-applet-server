package audit

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAuditTextScan(t *testing.T) {
	gone.Test(func(a *Audit) {
		pass, err := a.ScanText("哈哈哈")
		require.Nil(t, err)
		fmt.Printf("%v", pass)
	}, func(cemetery gone.Cemetery) error {
		goner.BasePriest(cemetery)
		cemetery.Bury(NewAudit())
		return nil
	})

}

func TestGmt(t *testing.T) {
	gmt := FormatGMT(time.Now())
	fmt.Println(gmt)
}
