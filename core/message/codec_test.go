package message

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMessage(t *testing.T) {
	var inputs = []string{
		`
{
    "id": 30349,
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10531,
    "seqId": 9,
    "sendId": "xxx",
    "channelId": 39,
    "content": {
        "reference": 30348,
        "text": "\n好的，我可以为您提供清明上河园游客组成比例和客流量的相关数据，此外，我还可以为您提供一些关于清明上河园的信息，比如清明上河园的历史背景、景点介绍等。",
        "type": "text"
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "url": "https://xxx/xxx",
        "type": "image"
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "type": "emoticon",
        "emoticonId": 1,
        "emoticonName": "表情名",
        "url": "https://xxx/xxx"
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "type": "emoticon-reply",
        "emoticonId": 1,
        "emoticonName": "大笑"
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "type": "file",
        "fileName": "吃瓜手册",
        "fileType": "pdf",
        "url": "https://xxx/xxx",
        "fileSize": 100
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "type": "channel-notice",
        "notice": "群公告"
    }
}
`,
		`
{
    "id": 30351,
    "sendAt": "2023-02-22T23:00:45.810524674+08:00",
    "createdAt": "2023-02-22T23:00:45.810524674+08:00",
    "userId": 10534,
    "sendId": "xxx",
    "seqId": 9,
    "channelId": 41,
    "content": {
        "reference": 30352,
        "type": "voice",
        "url": "https:///xxx/xxx",
        "duration": 31 
    }
}
`,
	}

	for _, i := range inputs {
		var m *Message
		err := json.Unmarshal([]byte(i), &m)
		if err != nil {
			t.Fatal(err)
			return
		}

		t.Logf("m: %v", m)

		j, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
			return
		}

		t.Logf("excepted: %s", i)
		t.Logf("actual: %s", j)

		if !compare([]byte(i), j) {
			t.Fatal("not equals between excepted and actual")
		}
	}
}

func compare(j1 []byte, j2 []byte) bool {
	var m1, m2 map[string]any
	err := json.Unmarshal(j1, &m1)
	if err != nil {
		return false
	}

	err = json.Unmarshal(j2, &m2)
	if err != nil {
		return false
	}

	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}

	return true
}
