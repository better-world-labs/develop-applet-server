package user

import (
	"fmt"
	"gitlab.openviewtech.com/gone/gone-lib/token"
)

func buildKey(moduleName string, businessName string, id interface{}) string {
	return fmt.Sprintf("%s#%s#%s", moduleName, businessName, id)
}

func buildSessionKey(userId int64, k string) string {
	key := "*"
	if k != "*" {
		key = token.ShortHash(k)
	}
	return buildKey("user", "k", fmt.Sprintf("%x:%s", userId, key))
}
