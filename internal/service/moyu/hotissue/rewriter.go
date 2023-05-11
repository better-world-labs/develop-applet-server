package hotissue

import (
	"fmt"
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"time"
)

type ContentProvider interface {
	func() string | string
}

type ContentRewriter struct {
	gone.Flag
	settings map[string]func() string
}

//go:gone
func NewContentRewriter() gone.Goner {
	r := ContentRewriter{}
	settings := map[string]func() string{
		"摸鱼八卦公会": func() string {
			return fmt.Sprintf("%d人同时在线，%d条热门讨论", utils.RandomInt(20, 100), utils.RandomInt(1000, 9000))
		},

		"闺蜜帮帮公会": func() string {
			return fmt.Sprintf("%d次情感问题讨论和分析", utils.RandomInt(10, 40))
		},

		"下班倒计时": func() string {
			later, _ := time.Parse("15:04:05", "18:00:00")
			nowTime, _ := utils.ParseClock(time.Now())
			return fmt.Sprintf("距离下班仅剩 %s", utils.FormatDuration(later.Sub(nowTime)))
		},
	}

	r.settings = settings
	return &r
}

func (r *ContentRewriter) Rewrite(issue *entity.HotIssue) {
	if setting, ok := r.settings[issue.Title]; ok {
		issue.Content = setting()
	}
}

func (r *ContentRewriter) RewriteList(issues []*entity.HotIssue) {
	for _, i := range issues {
		r.Rewrite(i)
	}
}
