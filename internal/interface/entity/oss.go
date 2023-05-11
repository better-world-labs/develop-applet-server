package entity

type OssToken struct {
	Host      string `json:"host"`
	AccessId  string `json:"accessId"`
	Signature string `json:"signature"`
	Policy    string `json:"policy"`
	Key       string `json:"key"`
}
