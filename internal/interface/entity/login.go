package entity

import "time"

const (
	UserIdKey    = "u-id"
	JwtKey       = "jwt"
	ClientIdKey  = "client-id"
	CsrfTokenKey = "X-CSRF-TOKEN"
)

type QrTokenWarp struct {
	QrToken string    `json:"qrToken"`
	Expired time.Time `json:"expired"`
}

type LoginRst struct {
	CsrfToken string      `json:"csrfToken"`
	User      *UserSimple `json:"user"`
}
