package miniapp

import (
	"bytes"
	"fmt"
	"github.com/gone-io/gone/goner/gin"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"io"
)

type AppRuntime struct {
	app             *entity.MiniApp
	currentFlow     int
	aiRuntime       *AIRuntime
	completeHandler func(outputs map[string]entity.MiniAppOutputCore)

	form    map[string]entity.MiniAppFormData
	outputs map[string]entity.MiniAppOutputCore

	outputStream chan *entity.MiniAppOutputStreamChunk
}

func NewAppRuntime(app *entity.MiniApp, aiRuntime *AIRuntime) *AppRuntime {
	return &AppRuntime{
		app:          app,
		outputStream: make(chan *entity.MiniAppOutputStreamChunk),
		aiRuntime:    aiRuntime,
		form:         make(map[string]entity.MiniAppFormData),
		outputs:      make(map[string]entity.MiniAppOutputCore),
	}
}

func (a *AppRuntime) initialize() error {
	for _, f := range *a.app.Form {
		if _, ok := a.form[f.GetId()]; ok {
			return errors.New(fmt.Sprintf("duplicate form %s\n", f.GetId()))
		}

		a.form[f.GetId()] = f
	}

	return nil
}

func (a *AppRuntime) parsePrompt(flow entity.MiniAppFlow) (string, error) {
	b := bytes.Buffer{}

	if flow.Prompt == nil {
		flow.Prompt = &entity.AppFlowPrompts{}
	}

	for _, pr := range *flow.Prompt {
		if tag, ok := pr.(*entity.AppFlowPromptTag); ok {
			switch tag.From {
			case entity.AppFlowPromptFromForm:
				if f, ok := a.form[tag.Character]; ok {
					b.WriteString(f.GetValue())
					continue
				}

			case entity.AppFlowPromptFromResult:
				if f, ok := a.outputs[tag.Character]; ok {
					b.WriteString(f.Content)
					continue
				}
			}

			return "", gin.NewParameterError("invalid tag")
		}

		if text, ok := pr.(*entity.AppFlowPromptText); ok {
			b.WriteString(text.Value)
		}
	}

	if b.Len() == 0 {
		return "", gin.NewParameterError("no prompt set")
	}

	return b.String(), nil
}

func (a *AppRuntime) runFlow(flow entity.MiniAppFlow) (TextStreamChunkReader, error) {
	prompt, err := a.parsePrompt(flow)
	if err != nil {
		return nil, err
	}

	reader, err := a.aiRuntime.createCompletion(flow.Type, prompt)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (a *AppRuntime) OnComplete(handler func(outputs map[string]entity.MiniAppOutputCore)) {
	a.completeHandler = handler
}

func (a *AppRuntime) Run() (*service.ChannelStreamTrunkReader[*entity.MiniAppOutputStreamChunk], error) {
	err := a.initialize()
	if err != nil {
		return nil, err
	}

	outputId := uuid.NewString()
	ch := make(chan *entity.MiniAppOutputStreamChunk)
	reader := service.NewChannelStreamTrunkReader[*entity.MiniAppOutputStreamChunk](ch)
	go func() {
		for _, flow := range a.app.Flow {
			textStreamReader, err := a.runFlow(flow)
			if err != nil {
				reader.SetInterrupt(err)
				return
			}

			var e error
			var output entity.MiniAppOutputCore
			for {
				textChunk, err := textStreamReader.Read()
				if err != nil {
					e = err
					break
				}

				if flow.OutputVisible {
					ch <- &entity.MiniAppOutputStreamChunk{
						OutputId: outputId,
						Flow:     flow.Id,
						Type:     entity.MiniAppOutputTypeText,
						Content:  string(textChunk),
					}
				}

				output.Flow = flow.Id
				output.OutputId = outputId
				output.Type = entity.MiniAppOutputTypeText
				output.Content = output.Content + string(textChunk) //TODO 优化
			}

			_ = textStreamReader.Close()
			a.outputs[flow.Id] = output
			if e != nil && e != io.EOF {
				reader.SetInterrupt(e)
				return
			}
		}

		a.completeHandler(a.outputs)
		close(ch)
	}()

	return reader, nil
}
