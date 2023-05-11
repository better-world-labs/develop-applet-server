package message

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"github.com/stretchr/testify/require"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"testing"
	"time"
)

func TestContextManager(t *testing.T) {
	gone.Test(func(m *contextManager) {
		err := m.WriteMessage(&message.Message{
			Header: message.Header{

				Id:        12,
				SendId:    "xxxx",
				CreatedAt: time.Now(),
				SendAt:    time.Now(),
				UserId:    10034,
				ChannelId: 40,
				SeqId:     1,
			},
			Content: message.NewTextContent(12, nil, "在这嗨"),
		})

		require.Nil(t, err)
		context, err := m.GetContext(40)
		require.Nil(t, err)

		for _, c := range context {
			require.True(t, c.Content != nil)
		}

	}, func(cemetery gone.Cemetery) error {
		_ = goner.BasePriest(cemetery)
		_ = goner.RedisPriest(cemetery)
		cemetery.Bury(NewContextManager())
		return nil
	})
}
