package entity

import "regexp"

type Uid string

var uidCheckReg = regexp.MustCompile(`^[0-9a-f]{6,16}$`)

func (u Uid) Check() bool {
	return uidCheckReg.Match([]byte(u))
}
