package notify_sender

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type NotifyMapping struct {
	NotifyType    entity.NotifyMessageType
	NotifyContent string
}

var notifies = map[string]NotifyMapping{
	entity.PointsTypeNewRegister: {
		NotifyType:    entity.NotifyMessageTypeReward,
		NotifyContent: "嗨新伙伴，欢迎加入我们！送你 %d 积分，和我们开始创造吧！",
	},
	entity.PointsTypeInvite: {
		NotifyType:    entity.NotifyMessageTypeInvited,
		NotifyContent: "哇塞！邀请好友成功，恭喜您获得 %d 积分",
	},
	entity.PointsTypeBeInvited: {
		NotifyType:    entity.NotifyMessageTypeInvited,
		NotifyContent: "新朋友，欢迎您！您被好友成功邀请，送你 %d 积分大礼包",
	},
	entity.PointsTypePointsRecharge: {
		NotifyType:    entity.NotifyMessageTypePointsRecharge,
		NotifyContent: "您充值的%d积分已到账~",
	},
	entity.PointsTypeAppDuplicated: {
		NotifyType:    entity.NotifyMessageDuplicateApp,
		NotifyContent: "您的小程序被一键同款啦，%d 积分已入账！",
	},
	entity.PointsTypeAppUsed: {
		NotifyType:    entity.NotifyMessageAppUsed,
		NotifyContent: "您的小程序获得了用户的认可，到账 %d 积分，以资鼓励！",
	},
}

func getNotifyMapping(_type string) (NotifyMapping, bool) {
	mapping, ok := notifies[_type]
	return mapping, ok
}
