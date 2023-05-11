package miniapp

import "bytes"

type AppCreatedNotifyTemplate struct {
	Env             string
	AppId           string
	Name            string
	Description     string
	DuplicateFrom   string
	CreatedNickname string
}

func (c AppCreatedNotifyTemplate) String() string {
	b := bytes.Buffer{}
	b.WriteString("数据统计: 新增小程序")
	b.WriteString("(")
	b.WriteString(c.Env)
	b.WriteString(" 环境)\n")
	b.WriteString("\n  - 小程序ID: ")
	b.WriteString(c.AppId)
	b.WriteString("\n  - 小程序名称: ")
	b.WriteString(c.Name)
	b.WriteString("\n  - 小程序描述: ")
	b.WriteString(c.Description)
	if len(c.DuplicateFrom) > 0 {
		b.WriteString("\n  - 使用模板: ")
		b.WriteString(c.DuplicateFrom)
	}
	b.WriteString("\n  - 创建用户: ")
	b.WriteString(c.CreatedNickname)
	return b.String()
}
