package recharge

import "testing"

func TestGenerateSerialNumber(t *testing.T) {
	number := GenerateSerialNumber()
	t.Logf("%s\n", number)
}
