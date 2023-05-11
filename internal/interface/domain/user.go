package domain

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

type (
	UserSettings struct {
		BossKey            string                     `json:"bossKey"`
		EndOffTime         string                     `json:"endOffTime"`
		AppearanceTheme    entity.AppearanceThemeEnum `json:"appearanceTheme"`
		SiteSettings       entity.SiteSettings        `json:"siteSettings"`
		MonthlySalary      int64                      `json:"monthlySalary"`
		MonthlyWorkingDays int64                      `json:"monthlyWorkingDays"`
	}

	EarlierThan struct {
		EarlierThan float64 `json:"earlierThan"`
	}

	AppearanceTheme struct {
		AppearanceTheme entity.AppearanceThemeEnum `json:"appearanceTheme"`
	}

	MoyuDetail struct {
		JoinDate             time.Time    `json:"joinDate"`
		TodayBrowseDuration  int64        `json:"todayBrowseDuration"`
		TotalBrowseDuration  int64        `json:"totalBrowseDuration"`
		MoreThan             float64      `json:"moreThan"`
		SecondSalary         float64      `json:"secondSalary"`
		LastReportBrowseTime time.Time    `json:"lastReportBrowseTime"`
		AccumulateMsgCnt     int64        `json:"accumulateMsgCnt"`
		User                 UserBaseInfo `json:"user"`
	}

	UserBaseInfo struct {
		Id       int64      `json:"id"`
		Nickname string     `json:"nickname"`
		Avatar   entity.Url `json:"avatar"`
	}

	WorkSettings struct {
		OffWorkTime        string `json:"offWorkTime"`
		MonthlySalary      int64  `json:"monthlySalary"`
		MonthlyWorkingDays int64  `json:"monthlyWorkingDays"`
	}

	RankInfo struct {
		Rank           int64        `json:"rank"`
		User           UserBaseInfo `json:"user"`
		BrowseDuration int64        `json:"browseDuration"`
	}
)
