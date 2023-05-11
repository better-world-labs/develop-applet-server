package entity

type StatisticMiniApp struct {
	AppId string `json:"-"`

	MiniAppStatisticInfo `xorm:"extends"`
}

type StatisticMiniAppOutput struct {
	OutputId string `json:"-"`

	MiniAppOutputStatisticInfo `xorm:"extends"`
}
