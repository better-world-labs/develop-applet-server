package trigger

import "github.com/gone-io/emitter"

type Sender struct {
	Sender emitter.Sender `gone:"gone-emitter"`
}

func (s *Sender) SendChannelTrigger(channelId int64) {

}

func (s *Sender) SendUserTrigger() {

}
