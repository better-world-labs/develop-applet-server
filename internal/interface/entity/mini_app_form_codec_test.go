package entity

import (
	"encoding/json"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	var forms = &MiniAppFormFields{
		NewMiniAppFormTextData("1", "姓名", "请输入姓名"),
		NewMiniAppFormSelectData("2", "性别", "请选择性别", "男\n女"),
	}

	b, err := json.Marshal(forms)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("Marshal result: %s\n", b)

	var unmarshal MiniAppFormFields
	err = json.Unmarshal(b, &unmarshal)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestEncodeDecode(t *testing.T) {
	var forms = []MiniAppFormData{
		NewMiniAppFormTextData("1", "姓名", "请输入姓名"),
		NewMiniAppFormSelectData("2", "性别", "请选择性别", "男\n女"),
	}

	b, err := EncodeMiniAppForms(forms)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("encode result: %s\n", b)
	decode, err := DecodeMiniAppForms(b)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("deocde result: %v\n", decode)
}
