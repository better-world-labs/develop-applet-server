package jssdk

import "testing"

func Test_createSignature(t *testing.T) {
	type args struct {
		nonceStr  string
		ticket    string
		timestamp string
		url       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				nonceStr:  "Wm3WZYTPz0wzccnW",
				timestamp: "1414587457",
				ticket:    "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
				url:       "http://mp.weixin.qq.com?params=value",
			},
			want: "0f9de62fce790f9a083d5c99e95740ceb90c27ed",
		},
		{
			name: "test#",
			args: args{
				nonceStr:  "Wm3WZYTPz0wzccnW",
				timestamp: "1414587457",
				ticket:    "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
				url:       "http://mp.weixin.qq.com?params=value#xxxx",
			},
			want: "0f9de62fce790f9a083d5c99e95740ceb90c27ed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSignature(tt.args.nonceStr, tt.args.ticket, tt.args.timestamp, tt.args.url); got != tt.want {
				t.Errorf("createSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractMainUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMainUrl(tt.args.url); got != tt.want {
				t.Errorf("extractMainUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
