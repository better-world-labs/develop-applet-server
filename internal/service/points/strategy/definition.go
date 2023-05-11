package strategy

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

var definition = map[string]entity.PointsDefinition{
	entity.PointsTypeNewRegister: {
		Type:        entity.PointsTypeNewRegister,
		Description: "注册奖励",
		Points:      +30,
	},
	entity.PointsTypeSignIn: {
		Type:        entity.PointsTypeSignIn,
		Description: "签到奖励",
		Points:      +5,
	},
	entity.PointsTypeInvite: {
		Type:        entity.PointsTypeInvite,
		Description: "邀请好友",
		Points:      +100,
	},
	entity.PointsTypeBeInvited: {
		Type:        entity.PointsTypeBeInvited,
		Description: "被好友邀请",
		Points:      +50,
	},
	entity.PointsTypeUsingApp: {
		Type:        entity.PointsTypeUsingApp,
		Description: "使用程序",
	},
	entity.PointsTypeAppUsed: {
		Type:        entity.PointsTypeAppUsed,
		Description: "程序被使用",
		Points:      +5,
	},
	entity.PointsTypeDuplicatingApp: {
		Type:        entity.PointsTypeDuplicatingApp,
		Description: "一键同款",
		Points:      -20,
	},
	entity.PointsTypeAppDuplicated: {
		Type:        entity.PointsTypeAppDuplicated,
		Description: "程序被一键同款",
		Points:      +20,
	},
	entity.PointsTypePointsRecharge: {
		Type:        entity.PointsTypePointsRecharge,
		Description: "积分充值",
	},
	entity.PointsTypeAppCreated: {
		Type:        entity.PointsTypeAppCreated,
		Description: "创建程序",
		Points:      +5,
	},
}

func getDefinition(_type string) (entity.PointsDefinition, bool) {
	def, ok := definition[_type]
	return def, ok
}
