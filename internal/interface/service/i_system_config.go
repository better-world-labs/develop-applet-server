package service

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type ISystemConfig interface {
	Put(key string, value any) error
	Get(Key string) (any, error)
}
