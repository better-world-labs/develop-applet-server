package entity

type ResignTemplate struct {
	Id      int64  `json:"id" xorm:"id pk autoincr"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
