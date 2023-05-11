package points

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/require"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"testing"
	"time"
)

func Priests(cemetery gone.Cemetery) error {
	_ = goner.BasePriest(cemetery)
	_ = goner.XormPriest(cemetery)
	cemetery.Bury(NewP())
	return nil
}

func TestCreatePoints(t *testing.T) {
	gone.Test(func(p *p) {
		err := p.create(&entity.Points{
			Points:      100,
			UserId:      10045,
			Description: "空头",
			CreatedAt:   time.Now(),
		})
		require.Nil(t, err)
	}, Priests)
}

func TestPagePointsFlow(t *testing.T) {
	gone.Test(func(p *p) {
		result, err := p.pageByUserId(page.Query{
			Page: 1,
			Size: 10,
		}, 10045)

		require.Nil(t, err)
		t.Log(result.Total, result.List)
	}, Priests)
}

func TestSumTotalPoints(t *testing.T) {
	gone.Test(func(p *p) {
		total, err := p.sumPointsByUserId(10045)
		require.Nil(t, err)
		t.Log(total)
	}, Priests)
}
