package entity

import (
	"errors"
	"fmt"
	"github.com/gone-io/gone/goner/gin"
	"time"
)

type MiniAppStatus uint8

const (
	MiniAppStatusUnPublished MiniAppStatus = 0
	MiniAppStatusPublished   MiniAppStatus = 1
)

type MiniAppCategory struct {
	Id   int64  `json:"id"`
	Text string `json:"text"`
}

type MiniAppBaseInfo struct {
	Id            int64         `json:"id"`
	Uuid          string        `json:"uuid"`
	DuplicateFrom string        `json:"duplicateFrom"`
	Name          string        `json:"name" binding:"required"`
	Description   string        `json:"description" binding:"required"`
	Category      int64         `json:"category" binding:"required"`
	Price         int           `xorm:"-" json:"price"`
	CreatedBy     int64         `json:"-"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Status        MiniAppStatus `json:"status"`
	Top           int           `json:"-"`
}

type MiniAppOutputStatisticInfo struct {
	LikeTimes    int `json:"likeTimes"`
	HateTimes    int `json:"hateTimes"`
	CommentTimes int `json:"commentTimes"`
}

type MiniAppStatisticInfo struct {
	RunTimes       int `json:"runTimes"`
	UseTimes       int `json:"useTimes"`
	LikeTimes      int `json:"likeTimes"`
	CommentTimes   int `json:"commentTimes"`
	CollectTimes   int `xorm:"-" json:"collectTimes"`
	ViewTimes      int `json:"viewTimes"`
	RecommendTimes int `json:"recommendTimes"`
	DegreeOfHeat   int `json:"degreeOfHeat"`
}

type MiniAppDetailDto struct {
	MiniApp
	MiniAppStatisticInfo

	SoldPoints int        `json:"soldPoints"`
	CreatedBy  UserSimple `json:"createdBy"`
}

type MiniApp struct {
	MiniAppBaseInfo `xorm:"extends"`

	Form *MiniAppFormFields `json:"form" xorm:"json"`
	Flow []MiniAppFlow      `json:"flow" xorm:"json"`
}

type MiniAppRunParam struct {
	Values []string `json:"values"`
	Open   bool     `json:"open"`
}

type MiniAppListDto struct {
	MiniAppBaseInfo
	MiniAppStatisticInfo

	SoldPoints int64            `json:"soldPoints"`
	Results    []*MiniAppOutput `json:"results"`
	Top        bool             `json:"top"`
	CreatedBy  UserSimple       `json:"createdBy"`
}

type (
	MiniAppOutputType string

	MiniAppOutputCore struct {
		OutputId string            `json:"id,omitempty"`
		Flow     string            `json:"flow,omitempty" xorm:"-"`
		Type     MiniAppOutputType `json:"type"`
		Content  string            `json:"content"`
	}

	TextStreamChunk []byte

	MiniAppOutputStreamChunk MiniAppOutputCore

	MiniAppOutput struct {
		MiniAppOutputCore `xorm:"extends"`

		Id        int64     `json:"-"`
		InputArgs []string  `json:"inputArgs"`
		Open      bool      `json:"-"`
		AppId     string    `json:"appId"`
		CreatedAt time.Time `json:"createdAt"`
		CreatedBy int64     `json:"-"`
	}

	MiniAppOutputDto struct {
		MiniAppOutput

		StatisticMiniAppOutput
		CreatedBy UserSimple `json:"createdBy"`
	}
)

type (
	MiniAppAiModelCategory struct {
		Id     int64             `json:"-"`
		Text   string            `json:"category"`
		Models []*MiniAppAiModel `json:"models" xorm:"-"`
	}

	MiniAppAiModel struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Category    int64  `json:"category,omitempty"`
		Available   bool   `json:"available"`
	}
)

const (
	MiniAppOutputTypeText  MiniAppOutputType = "text"
	MiniAppOutputTypeImage MiniAppOutputType = "image"
)

func (a *MiniAppBaseInfo) IsDuplicated() bool {
	return a.DuplicateFrom != ""
}

func (a *MiniAppBaseInfo) Cursor() int64 {
	return a.Id
}

func (a *MiniAppBaseInfo) IsTop() bool {
	return a.Top > 0
}

func (a *MiniAppOutput) Cursor() int64 {
	return a.Id
}

func (a *MiniApp) Input(args []string) error {
	if len(*a.Form) != len(args) {
		return errors.New("invalid args")
	}

	for i, field := range *a.Form {
		err := field.SetValue(args[i])
		if err != nil {
			return gin.NewParameterError(fmt.Sprintf("set value error on [%d] arg: %s\n", i, err.Error()))
		}
	}

	return nil
}

type MiniAppOutputStreamReader interface {
	Read() (MiniAppOutputCore, error)
	OnComplete(handler func(MiniAppOutputCore))
}
