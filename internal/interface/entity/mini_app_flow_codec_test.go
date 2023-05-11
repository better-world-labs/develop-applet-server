package entity

import (
	"encoding/json"
	"testing"
)

func TestMarshalUnmarshalAppFlow(t *testing.T) {
	var forms = &AppFlowPrompts{
		NewAppFlowPromptText("从"),
		NewAppFlowPromptTag(AppFlowPromptFromResult, "1"),
	}

	b, err := json.Marshal(forms)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("Marshal result: %s\n", b)

	var unmarshal AppFlowPrompts
	err = json.Unmarshal(b, &unmarshal)
	if err != nil {
		t.Fatal(err)
		return
	}
}
func TestAppFlowEncodeDecode(t *testing.T) {
	var prompts = []AppFlowPrompt{
		NewAppFlowPromptText("从"),
		NewAppFlowPromptTag(AppFlowPromptFromResult, "1"),
	}

	b, err := EncodeMiniAppFlowPrompts(prompts)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("encode result: %s\n", b)
	decode, err := DecodeMiniAppFlowPrompts(b)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("deocde result: %v\n", decode)
}
