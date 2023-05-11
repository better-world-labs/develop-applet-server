package message

// messageContentFactory 提供不同类型 Content 结构体
var messageContentFactory = map[ContentType]func() Content{
	ContentTypeText: func() Content {
		return &TextContent{}
	},
	ContentTypeImage: func() Content {
		return &ImageContent{}
	},
	ContentTypeEmoticon: func() Content {
		return &EmoticonContent{}
	},
	ContentTypeEmoticonReply: func() Content {
		return &EmoticonReplyContent{}
	},
	ContentTypeFile: func() Content {
		return &FileContent{}
	},
	ContentTypeChannelNotice: func() Content {
		return &ChannelNoticeContent{}
	},
	ContentTypeVoice: func() Content {
		return &VoiceContent{}
	},
}
