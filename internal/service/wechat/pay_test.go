package wechat

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/assert"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"os"
	"testing"
	"time"
)

func Test_paySvc_CreateOrderForNativeAPI(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Printf("workDir:%s", dir)

	gone.Test(func(s *paySvc) {
		const Attach = "my-test"
		const orderNo = "11011101fads11"
		res, err := s.CreateOrderForNativeAPI(&entity.WxPayOrderRequest{
			OrderNo:     orderNo,
			Description: "测试商品信息",
			TimeExpire:  time.Now().Add(10 * time.Minute),
			Total:       1,
			Openid:      "omxLu5AUGP-JJ6GlrbBMBuQyEJL4",
			Attach:      Attach,
			NotifyUrl:   "https://foods-dev.openviewtech.com/pays/notify",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		t.Logf("codeUrl=%s\n", *res)

		trans, _, err := s.QueryOrderInfoByOrderNo(orderNo)
		assert.Nil(t, err)
		assert.NotEmpty(t, trans.TradeState)

		//trans2, err := s.QueryOrderInfoByWxTransactionId(*trans.TransactionId)
		//assert.Nil(t, err)
		//assert.Equal(t, *trans2.TransactionId, *trans.TransactionId)
	}, func(cemetery gone.Cemetery) error {
		_ = goner.BasePriest(cemetery)
		cemetery.Bury(NewPaySvc())
		return nil
	})
}
func Test_paySvc_CreateOrderForJSAPI(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Printf("workDir:%s", dir)

	gone.Test(func(s *paySvc) {
		const Attach = "my-test"
		const orderNo = "1101110s111"
		res, err := s.CreateOrderForJSAPI(&entity.WxPayOrderRequest{
			OrderNo:     orderNo,
			Description: "测试商品信息",
			TimeExpire:  time.Now().Add(10 * time.Minute),
			Total:       10,
			Openid:      "omxLu5AUGP-JJ6GlrbBMBuQyEJL4",
			Attach:      Attach,
			NotifyUrl:   "https://foods-dev.openviewtech.com/pays/notify",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, res.PrepayId)

		trans, _, err := s.QueryOrderInfoByOrderNo(orderNo)
		assert.Nil(t, err)
		assert.NotEmpty(t, trans.TradeState)

		//trans2, err := s.QueryOrderInfoByWxTransactionId(*trans.TransactionId)
		//assert.Nil(t, err)
		//assert.Equal(t, *trans2.TransactionId, *trans.TransactionId)
	}, func(cemetery gone.Cemetery) error {
		_ = goner.BasePriest(cemetery)
		cemetery.Bury(NewPaySvc())
		return nil
	})
}

func TestQueryOrder(t *testing.T) {
	dir, _ := os.Getwd()
	fmt.Printf("workDir:%s", dir)

	gone.Test(func(s *paySvc) {
		trans, _, err := s.QueryOrderInfoByOrderNo("e35a046d5c0047339ab2f871ae5d8d1e")
		assert.Nil(t, err)
		assert.NotEmpty(t, trans.TradeState)

	}, func(cemetery gone.Cemetery) error {
		_ = goner.BasePriest(cemetery)
		cemetery.Bury(NewPaySvc())
		return nil
	})
}
