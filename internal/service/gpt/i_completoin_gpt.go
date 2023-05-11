package gpt

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type (
	ICompletionGpt interface {
		CreateChatCompletionStream(request ChatCompletionRequest) (*ChatCompletionStreamReader, error)
	}
)
