package message

// TextContent 文本
type TextContent struct {
	ContentHeader

	Text string `json:"text"`
}

func NewTextContent(referenceId int64, reply []Reply, text string) *TextContent {
	return &TextContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeText, reply),
		Text:          text,
	}
}

// ImageContent 图片
type ImageContent struct {
	ContentHeader

	Url string `json:"url"`
}

func NewImageContent(referenceId int64, reply []Reply, url string) *ImageContent {
	return &ImageContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeImage, reply),
		Url:           url,
	}
}

// EmoticonContent 表情
type EmoticonContent struct {
	ContentHeader

	EmoticonId   int64  `json:"emoticonId"`
	EmoticonName string `json:"emoticonName"`
	Url          string `json:"url"`
}

func NewEmoticonContent(referenceId int64, reply []Reply, emotionId int64, emotionName string, url string) *EmoticonContent {
	return &EmoticonContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeEmoticon, reply),
		EmoticonId:    emotionId,
		EmoticonName:  emotionName,
		Url:           url,
	}
}

// FileContent 文件
type FileContent struct {
	ContentHeader

	FileType string `json:"fileType"`
	FileName string `json:"fileName"`
	Url      string `json:"url"`
	FileSize int    `json:"fileSize"`
}

func NewFileContent(referenceId int64, reply []Reply, fileType, fileName string, fileSize int, url string) *FileContent {
	return &FileContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeFile, reply),
		FileType:      fileType,
		FileName:      fileName,
		FileSize:      fileSize,
		Url:           url,
	}
}

// EmoticonReplyContent 表情回复
type EmoticonReplyContent struct {
	ContentHeader

	EmoticonId   int64  `json:"emoticonId"`
	EmoticonName string `json:"emoticonName"`
}

func NewEmoticonReplyContent(referenceId int64, reply []Reply, emoticonId int64, emoticonName string) *EmoticonReplyContent {
	return &EmoticonReplyContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeEmoticonReply, reply),
		EmoticonId:    emoticonId,
		EmoticonName:  emoticonName,
	}
}

// ChannelNoticeContent 频道公告
type ChannelNoticeContent struct {
	ContentHeader

	Notice string `json:"notice"`
}

func NewChannelNoticeContent(referenceId int64, reply []Reply, notice string) *ChannelNoticeContent {
	return &ChannelNoticeContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeChannelNotice, reply),
		Notice:        notice,
	}
}

// VoiceContent 语音
type VoiceContent struct {
	ContentHeader

	Url      string `json:"url"`
	Duration int    `json:"duration"`
}

func NewVoiceContent(referenceId int64, reply []Reply, url string, duration int) *VoiceContent {
	return &VoiceContent{
		ContentHeader: *NewContentHeader(referenceId, ContentTypeVoice, reply),
		Url:           url,
		Duration:      duration,
	}
}
