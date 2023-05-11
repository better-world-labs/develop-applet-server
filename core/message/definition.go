package message

import (
	"encoding/json"
	"time"
)

const (
	ContentTypeText          ContentType = "text"
	ContentTypeImage         ContentType = "image"
	ContentTypeEmoticon      ContentType = "emoticon"
	ContentTypeFile          ContentType = "file"
	ContentTypeEmoticonReply ContentType = "emoticon-reply"
	ContentTypeChannelNotice ContentType = "channel-notice"
	ContentTypeVoice         ContentType = "voice"
)

type (
	ContentType string

	// Message 消息在业务流转的形式定义
	Message struct {
		Header `xorm:"extends"`

		Content Content `xorm:"json" json:"content"`
	}

	// Header 消息头部
	Header struct {
		Id        int64     `json:"id"`
		SendId    string    `json:"sendId"`
		CreatedAt time.Time `json:"createdAt"`
		SendAt    time.Time `json:"sendAt"`
		SeqId     int64     `json:"seqId"`
		UserId    int64     `json:"userId"`
		ChannelId int64     `json:"channelId"`
	}

	// Content 消息内容抽象
	Content interface {
		GetReference() int64
		GetType() ContentType
		GetReply() []Reply
	}

	ContentHeader struct {
		Reference int64       `json:"reference,omitempty"`
		Type      ContentType `json:"type"`
		Reply     []Reply     `json:"reply,omitempty"`
	}

	Reply struct {
		UserId     int64 `json:"userId"`
		EmoticonId int   `json:"emoticonId"`
	}
)

func NewContentHeader(reference int64, contentType ContentType, reply []Reply) *ContentHeader {
	return &ContentHeader{
		Reference: reference,
		Type:      contentType,
		Reply:     reply,
	}
}

func (m ContentHeader) GetReference() int64 {
	return m.Reference
}

func (m ContentHeader) GetReply() []Reply {
	return m.Reply
}

func (m ContentHeader) GetType() ContentType {
	return m.Type
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var s struct {
		Header

		Content *json.RawMessage `json:"content"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if m == nil {
		*m = Message{}
	}

	m.Header = s.Header

	if s.Content != nil {
		content, _ := DecodeContent(*s.Content)
		if err != nil {
			return err
		}

		m.Content = content
	}

	return nil
}
