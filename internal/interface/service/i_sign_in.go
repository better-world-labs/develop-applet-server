package service

type ISignIn interface {
	GetSignInStatus(userId int64) (bool, error)
	SignIn(userId int64) error
}
