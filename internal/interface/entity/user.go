package entity

import "time"

type AppearanceThemeEnum string

const (
	AppearanceThemeBright AppearanceThemeEnum = "bright"
	AppearanceThemeDark   AppearanceThemeEnum = "dark"
)

type User struct {
	Id                  int64      `json:"id" xorm:"id pk autoincr"`
	WxOpenId            string     `json:"wxOpenId"`
	WxUnionId           string     `json:"wxUnionId"`
	Nickname            string     `json:"nickname"`
	Avatar              Url        `json:"avatar"`
	Online              bool       `json:"online"`
	InvitedBy           *int64     `json:"invitedBy"`
	FromApp             string     `json:"-"`
	Source              string     `json:"-"`
	ConnectTime         *time.Time `json:"connectTime"`
	LoginAt             *time.Time `json:"loginAt"`
	LastLoginAt         *time.Time `json:"lastLoginAt"`
	TotalAccessDuration int64      `json:"totalAccessDuration"`
	TotalBrowseDuration int64      `json:"totalBrowseDuration"`
	CreatedAt           time.Time  `json:"createdAt"`
}

func (u User) RegisteredDays() int {
	return int(time.Now().Sub(u.CreatedAt).Hours() / 24)
}

func (u User) IsFirstLogin() bool {
	return u.LastLoginAt == nil
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Avatar   Url    `json:"avatar"`
	Sex      string `json:"sex"`
}

type UserInfo struct {
	UserSimple
	LoginAt     *time.Time `json:"loginAt"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	InvitedBy   int64      `json:"invitedBy"`
	Points      int64      `json:"points"`
}

type UserSimple struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   Url    `json:"avatar"`
}

type UserSettings struct {
	Id                 int64               `json:"id"`
	UserId             int64               `json:"userId"`
	EndOffTime         string              `json:"endOffTime"`
	BossKey            string              `json:"bossKey"`
	AppearanceTheme    AppearanceThemeEnum `json:"appearanceTheme"`
	SiteSettings       string              `json:"siteSettings"`
	MonthlySalary      int64               `json:"monthlySalary"`
	MonthlyWorkingDays int64               `json:"monthlyWorkingDays"`
}

type SiteSettings struct {
	Type        SiteSettingsType `json:"type"`
	CustomTitle string           `json:"customTitle"`
	CustomIcon  string           `json:"customIcon"`
}

type SiteSettingsType string

const (
	SiteSettingsDefault SiteSettingsType = "default"
	SiteSettingsOffice  SiteSettingsType = "office"
	SiteSettingsCustom  SiteSettingsType = "custom"
)

type Identity struct {
	Nickname string `json:"nickname"`
	Avatar   Url    `json:"avatar"`
}

type ComponentNameEnum string

const (
	AppearanceThemeComponent ComponentNameEnum = "appearanceTheme"
	SiteSettingsComponent    ComponentNameEnum = "siteSettings"
	MonthlySalary            ComponentNameEnum = "monthlySalary"
	OffWorkTime              ComponentNameEnum = "offWorkTime"
	MonthlyWorkingDays       ComponentNameEnum = "monthlyWorkingDays"
)

func (cn *ComponentNameEnum) ToQueryCol() string {
	var col string
	switch *cn {
	case AppearanceThemeComponent:
		col = "appearance_theme"
	case SiteSettingsComponent:
		col = "site_settings"
	case MonthlySalary:
		col = "monthly_salary"
	case OffWorkTime:
		col = "end_off_time"
	case MonthlyWorkingDays:
		col = "monthly_working_days"
	}

	return col
}

func (cn *ComponentNameEnum) ToComponentElement(settings UserSettings) ComponentElement {
	var element ComponentElement
	switch *cn {
	case AppearanceThemeComponent:
		element.ComponentName = AppearanceThemeComponent
		element.ComponentSettings = string(settings.AppearanceTheme)
	case SiteSettingsComponent:
		element.ComponentName = SiteSettingsComponent
		element.ComponentSettings = settings.SiteSettings
	}

	return element
}

func (cn *ComponentNameEnum) ValidComponent() bool {
	return *cn == AppearanceThemeComponent || *cn == SiteSettingsComponent
}

type ComponentElement struct {
	ComponentName     ComponentNameEnum `json:"componentName"`
	ComponentSettings interface{}       `json:"componentSettings"`
}

type UserAccessRecord struct {
	Id        int64      `json:"id" xorm:"id pk autoincr"`
	UserId    int64      `json:"userId"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}

type UserBrowserRecord struct {
	Id             int64     `json:"id" xorm:"id pk autoincr"`
	UserId         int64     `json:"userId"`
	BrowseDate     string    `json:"browseDate"`
	BrowseDuration int64     `json:"browseDuration"`
	LastBrowseTime time.Time `json:"lastBrowseTime"`
}

type TimeQuantum struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type StatPeriod int64

const (
	Total StatPeriod = iota
	Daily
)
