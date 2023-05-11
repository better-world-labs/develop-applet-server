package entity

type Emoticon struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      Url    `json:"url"`
	Keywords string `json:"keywords"`
}
