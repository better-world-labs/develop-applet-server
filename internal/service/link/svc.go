package link

import (
	"encoding/base32"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/redis"
	"github.com/google/uuid"
	"time"
)

const (
	DefaultTTL = time.Hour * 24
)

type svc struct {
	gone.Flag
	Cache redis.Cache `gone:"gone-redis-cache"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s *svc) Create(origin string) (string, error) {
	code := genericLinkCode()
	return code, s.Cache.Put(key(code), origin, DefaultTTL)
}

func (s *svc) GetOrigin(linkCode string) (string, bool, error) {
	var res *string
	err := s.Cache.Get(key(linkCode), &res)
	if err != nil {
		if err == redis.ErrNil {
			return "", false, nil
		}

		return "", false, err
	}

	return *res, res != nil, nil
}

func genericLinkCode() string {
	u := uuid.NewString()
	str := []byte(u)[:6]
	return base32.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", str, time.Now().Unix())))
}

func key(key string) string {
	return fmt.Sprintf("LINK-%s", key)
}
