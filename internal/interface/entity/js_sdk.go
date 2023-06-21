package entity

type JSSDKSignature struct {
	NonceStr  string `json:"nonceStr"`
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	AppId     string `json:"appId"`
}
