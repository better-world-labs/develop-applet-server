package stat

type hotMsg struct {
	MsgId      int64
	ReplyCount int
}

//go:generate sh -c "mockgen -package=$GOPACKAGE -self_package=moyu-server/internal/service/$GOPACKAGE  -source=$GOFILE|gone mock -o persistence_mock_test.go"
type iPersistence interface {
	listTopReplyMsg(top int, channelId int64) ([]*hotMsg, error)
}
