package gptconversation

import (
	"bytes"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/gpt"
	"io"
	"time"
)

const DefaultChatContext = 5

type svc struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`

	points service.IPointStrategy `gone:"*"`
	p      iPersistence           `gone:"*"`
	gpt    gpt.ICompletionGpt     `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) createContext(message entity.GptChatMessage) ([]gpt.ChatCompletionMessage, error) {
	var res []gpt.ChatCompletionMessage

	result, err := s.p.pageByUserId(page.CreateStreamQuery(DefaultChatContext, ""), message.UserId)
	if err != nil {
		return nil, err
	}

	messages := result.GetList()
	messages = append([]*entity.GptChatMessage{&message}, messages...)

	linq.From(result.GetList()).Sort(func(i, j interface{}) bool {
		return i.(*entity.GptChatMessage).Id < j.(*entity.GptChatMessage).Id
	}).Select(func(i interface{}) interface{} {
		message := i.(*entity.GptChatMessage)
		return gpt.ChatCompletionMessage{
			Role:    string(message.Role),
			Content: message.Content,
		}
	}).ToSlice(&res)

	res = append(res)
	return res, nil
}

func (s svc) SendMessage(userId int64, content string) (*service.ChannelStreamTrunkReader[entity.GptChatMessage], error) {
	message := entity.GptChatMessage{
		Role:      entity.GPTRoleUser,
		Content:   content,
		UserId:    userId,
		CreatedAt: time.Now(),
	}

	if err := s.Transaction(func(session xorm.Interface) error {
		if err := s.p.create(&message); err != nil {
			return err
		}

		_, err := s.points.ApplyPoints(userId, entity.StrategyArgGptConversation{})
		return err
	}); err != nil {
		return nil, err
	}

	context, err := s.createContext(message)
	if err != nil {
		return nil, err
	}
	stream, err := s.gpt.CreateChatCompletionStream(gpt.ChatCompletionRequest{
		Model:            "gpt-3.5-turbo",
		Temperature:      0.9,
		MaxTokens:        2048,
		TopP:             0.9,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Messages:         context,
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan entity.GptChatMessage)
	reader := service.NewChannelStreamTrunkReader(ch)
	go func() {
		var replyContent = bytes.Buffer{}

		for {
			r, err := stream.Read()
			if err != nil {
				reader.SetInterrupt(err)
				if err == io.EOF {
					stream.Close()
				}

				break
			}

			replyChunk := entity.GptChatMessage{
				Role:    entity.GPTRoleAssistant,
				UserId:  userId,
				Content: r.Choices[0].Delta.Content,
			}
			ch <- replyChunk
			replyContent.WriteString(replyChunk.Content)
		}

		if err := s.p.create(&entity.GptChatMessage{
			Role:      entity.GPTRoleAssistant,
			UserId:    userId,
			Content:   replyContent.String(),
			CreatedAt: time.Now(),
		}); err != nil {
			//TODO
		}
	}()

	return reader, nil
}

func (s svc) ListMessages(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.GptChatMessage], error) {
	isEmpty, err := s.p.isEmptyConversation(userId)
	if err != nil {
		return nil, err
	}

	if isEmpty {
		err := s.initGreetMessage(userId)
		if err != nil {
			return nil, err
		}
	}

	return s.p.pageByUserId(query, userId)
}

func (s svc) initGreetMessage(userId int64) error {
	return s.p.create(&entity.GptChatMessage{
		Role:      entity.GPTRoleAssistant,
		UserId:    userId,
		Content:   "你好呀，请试试向我提问，例如：现在请你扮演一位严格对待员工，非要紧事情不允许员工请假的企业领导。而我作为员工想要在工作日向你请假，虽然我请假的目的只是为了出去玩或者在家休息。请从帮我写一段请假理由，要求是当你看到这个请假理由后会欣然同意我的请假申请，字数在100个汉字以内。",
		CreatedAt: time.Now(),
	})
}
