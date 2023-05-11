package stat

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	message3 "gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service/mock"
	"testing"
)

func Test_messageStat_HotMessageTop(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPersistence := NewMockiPersistence(ctrl)

	mockRecord := mock.NewMockIMessageRecord(ctrl)

	stat := messageStat{
		p:        mockPersistence,
		iMessage: mockRecord,
	}

	mockPersistence.EXPECT().
		listTopReplyMsg(10, int64(25)).Return([]*hotMsg{
		{MsgId: 10, ReplyCount: 1},
		{MsgId: 11, ReplyCount: 2},
		{MsgId: 12, ReplyCount: 3},
		{MsgId: 13, ReplyCount: 4},
		{MsgId: 14, ReplyCount: 5},
	}, nil)

	mockRecord.EXPECT().
		GetRecordsMap([]int64{
			10, 11, 12, 13, 14,
		}).
		Return(map[int64]*message3.Message{
			10: {
				Header: message3.Header{
					Id:        10,
					UserId:    10,
					ChannelId: 25,
				},
				Content: message3.TextContent{
					ContentHeader: *message3.NewContentHeader(0, message3.ContentTypeText, nil),
					Text:          "this is a test",
				},
			},
			11: {
				Header: message3.Header{
					Id:        10,
					UserId:    10,
					ChannelId: 25,
				},
				Content: &message3.EmoticonContent{
					ContentHeader: *message3.NewContentHeader(0, message3.ContentTypeEmoticon, nil),
					EmoticonId:    10,
					EmoticonName:  "test",
				},
			},
			12: {
				Header: message3.Header{
					Id:        10,
					UserId:    10,
					ChannelId: 25,
				},
				Content: &message3.TextContent{
					ContentHeader: *message3.NewContentHeader(0, message3.ContentTypeText, nil),
					Text:          "this is a test",
				},
			},
			13: {
				Header: message3.Header{
					Id:        10,
					UserId:    10,
					ChannelId: 25,
				},
				Content: &message3.TextContent{
					ContentHeader: *message3.NewContentHeader(0, message3.ContentTypeText, nil),
					Text:          "this is a test",
				},
			},
			14: {
				Header: message3.Header{
					Id:        10,
					UserId:    10,
					ChannelId: 25,
				},
				Content: &message3.TextContent{
					ContentHeader: *message3.NewContentHeader(0, message3.ContentTypeText, nil),
					Text:          "this is a test",
				},
			},
		}, nil)

	out, err := stat.HotMessageTop(25, 10)
	assert.Nil(t, err)

	assert.Equal(t, len(out), 5)

	assert.Equal(t, out[0].MsgId, int64(10))
	assert.Equal(t, out[0].Content.GetType(), message3.ContentType("text"))

	message, ok := out[1].Content.(*message3.EmoticonContent)
	assert.True(t, ok)
	assert.Equal(t, message.EmoticonId, int64(10))
}

func encode(content message3.Content) []byte {
	b, _ := json.Marshal(content)
	return b
}
