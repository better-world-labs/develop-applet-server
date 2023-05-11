package sign_in

import "time"

type SignInDaily struct {
	Id     int64
	UserId int64
	Date   time.Time
}
