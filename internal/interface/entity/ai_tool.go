package entity

type AiTool struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    int    `json:"category"`
	Target      string `json:"target"`
}

type AiToolCategory struct {
	Category    int    `json:"category"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}
