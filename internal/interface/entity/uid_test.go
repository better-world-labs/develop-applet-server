package entity

import "testing"

func TestUid_Check(t *testing.T) {
	tests := []struct {
		name string
		u    Uid
		want bool
	}{
		{
			name: "correct",
			u:    "3f8e34f4",
			want: true,
		},
		{
			name: "incorrect",
			u:    "3f8e34f4x",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Check(); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
