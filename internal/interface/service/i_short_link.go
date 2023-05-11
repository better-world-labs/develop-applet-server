package service

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IShortLink interface {
	Create(origin string) (string, error)

	GetOrigin(linkCode string) (string, bool, error)
}
