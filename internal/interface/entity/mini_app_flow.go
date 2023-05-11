package entity

type (
	AppFlowPrompts []AppFlowPrompt
	MiniAppFlow    struct {
		Id            string          `json:"id"`
		Type          string          `json:"type"`
		OutputVisible bool            `json:"outputVisible"`
		Prompt        *AppFlowPrompts `json:"prompt"`
	}
)

func (m AppFlowPrompts) MarshalJSON() ([]byte, error) {
	return EncodeMiniAppFlowPrompts(m)
}

func (m *AppFlowPrompts) UnmarshalJSON(b []byte) error {
	prompt, err := DecodeMiniAppFlowPrompts(b)
	if err != nil {
		return err
	}

	*m = prompt
	return nil
}

type (
	AppFlowPromptType string
	AppFlowFromType   string
)

const (
	AppFlowPromptTypeText   AppFlowPromptType = "text"
	AppFlowPromptTypeTag    AppFlowPromptType = "tag"
	AppFlowPromptFromResult AppFlowFromType   = "result"
	AppFlowPromptFromForm   AppFlowFromType   = "form"
)

type AppFlowPrompt interface {
	GetType() AppFlowPromptType
}

type AppFlowPromptText struct {
	Value string `json:"value"`
}

func NewAppFlowPromptText(value string) *AppFlowPromptText {
	return &AppFlowPromptText{Value: value}
}

func (a AppFlowPromptText) GetType() AppFlowPromptType {
	return AppFlowPromptTypeText
}

type AppFlowPromptTag struct {
	From      AppFlowFromType `json:"from"`
	Character string          `json:"character"`
}

func NewAppFlowPromptTag(from AppFlowFromType, character string) *AppFlowPromptTag {
	return &AppFlowPromptTag{
		From:      from,
		Character: character,
	}
}

func (a AppFlowPromptTag) GetType() AppFlowPromptType {
	return AppFlowPromptTypeTag
}
