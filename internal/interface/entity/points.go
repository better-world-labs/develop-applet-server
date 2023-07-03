package entity

import (
	"math"
	"time"
)

const (
	PointsTypeNewRegister     = "new-register"
	PointsTypeSignIn          = "sign-in"
	PointsTypeInvite          = "invite"
	PointsTypeBeInvited       = "be-invited"
	PointsTypeUsingApp        = "using-app"
	PointsTypeAppUsed         = "app-used"
	PointsTypeAppCreated      = "app-created"
	PointsTypeGptConversation = "apt-conversation"
	PointsTypeDuplicatingApp  = "duplicating-app"
	PointsTypeAppDuplicated   = "app-duplicated"
	PointsTypePointsRecharge  = "points-recharge"
)

type PointsDefinition struct {
	Type        string
	Description string
	Points      int
	Strategy    any
}

type PointsRanking struct {
	UserId int64      `json:"-"`
	User   UserSimple `json:"user"`
	Points int64      `json:"points"`
}

type Points struct {
	Id          int64     `json:"id"`
	Points      int64     `json:"points"`
	Type        string    `json:"type"`
	UserId      int64     `json:"userId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (p Points) AbsPoints() int {
	return int(math.Abs(float64(p.Points)))
}
