package message

import (
	"encoding/json"
	"github.com/gone-io/gone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
)

func getMsgReplyLength(t *testing.T, p *persistence, msgId int64) int {
	msg, b, err := p.getById(msgId)
	assert.Nil(t, err)
	assert.True(t, b)

	assert.Nil(t, err)
	reply := msg.Content.GetReply()
	return len(reply)
}

func Test_insertReplyForMessage(t *testing.T) {
	gone.Test(func(p *persistence) {
		const msgId = 1399
		l1 := getMsgReplyLength(t, p, msgId)
		p.insertReplyForMessage(msgId, 10, 10034)
		l2 := getMsgReplyLength(t, p, msgId)
		assert.Equal(t, l2, l1+1)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}

func Test_persistence_GetRecordsSummary(t *testing.T) {
	gone.Test(func(p *persistence) {
		message, err := p.GetRecordsSummary(1400, 41)
		require.Nil(t, err)

		b, err := json.Marshal(message)
		require.Nil(t, err)
		t.Logf("message: %s", b)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}
func Test_persistence_ListByIdsMap(t *testing.T) {
	gone.Test(func(p *persistence) {
		message, err := p.listByIdsMap([]int64{1418})
		require.Nil(t, err)

		b, err := json.Marshal(message)
		require.Nil(t, err)
		t.Logf("message: %s", b)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}
func Test_persistence_ListHistoryMessageAfter(t *testing.T) {
	gone.Test(func(p *persistence) {
		message, err := p.listHistoryMessageAfter(41, 1418, 100)
		require.Nil(t, err)

		b, err := json.Marshal(message)
		require.Nil(t, err)
		t.Logf("message: %s", b)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}
func Test_persistence_ListHistoryMessageBefore(t *testing.T) {
	gone.Test(func(p *persistence) {
		message, err := p.listHistoryMessageBefore(41, 1418, 100)
		require.Nil(t, err)

		b, err := json.Marshal(message)
		require.Nil(t, err)
		t.Logf("message: %s", b)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}

func Test_persistence_GetLastMessageByChannelId(t *testing.T) {
	gone.Test(func(p *persistence) {
		message, err := p.GetLastMessageByChannelId(41)
		require.Nil(t, err)

		b, err := json.Marshal(message)
		require.Nil(t, err)
		t.Logf("message: %s", b)
	}, func(cemetery gone.Cemetery) error {

		//引入测试用的mysqlPriest
		_ = test.MysqlPriest(cemetery)

		//注入当前需要测试的持久类
		cemetery.Bury(NewPersistence())
		return nil
	})
}
