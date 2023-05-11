package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/redis"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
)

const (
	KeyPrefix = "Context_Channel_"
)

// LuaScriptPush 入队 ARGV[1] = 值， ARGV[2] = 队列最大长度
const LuaScriptPush = `
if redis.call('LPUSH', KEYS[1], ARGV[1]) then
   return redis.call('LTRIM', KEYS[1], 0, ARGV[2] - 1)
end
   return nil
`

type contextManager struct {
	gone.Goner

	Redis            redis.Pool `gone:"gone-redis-pool"`
	maxContextLength int        `gone:"config,cache.max-context-length"`
}

//go:gone
func NewContextManager() gone.Goner {
	return &contextManager{}
}

func (c contextManager) WriteMessage(record *message.Message) error {
	conn := c.Redis.Get()
	defer conn.Close()

	b, err := json.Marshal(record)
	if err != nil {
		return err
	}

	reply, err := conn.Do("EVAL", LuaScriptPush, 1, c.buildKey(record.ChannelId), b, c.maxContextLength)
	if err != nil {
		return err
	}

	if reply != nil && reply.(string) == "OK" {
		return nil
	}

	return errors.New("LPUSH failed")
}

func (c contextManager) GetContext(channelId int64) ([]*message.Message, error) {
	conn := c.Redis.Get()
	defer conn.Close()

	reply, err := conn.Do("LRANGE", c.buildKey(channelId), 0, -1)
	if err != nil {
		return nil, err
	}

	bs := reply.([]any)
	messages := make([]*message.Message, 0, len(bs))
	for _, b := range bs {
		var message message.Message
		err := json.Unmarshal(b.([]byte), &message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

func (c contextManager) buildKey(channelId int64) string {
	return fmt.Sprintf("%s_%d", KeyPrefix, channelId)
}
