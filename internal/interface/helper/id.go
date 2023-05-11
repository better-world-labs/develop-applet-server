package helper

import (
	"github.com/google/uuid"
	"strings"
)

func GenerateSerialNumber() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))
}
