package notice

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {

	ids, message, err := parseMentionedUserIds("@[1001] @[1002] 哈 哈哈")
	if err != nil {
		panic(err)
	}

	fmt.Println(ids)
	fmt.Println(message)
}
