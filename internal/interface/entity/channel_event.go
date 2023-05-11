package entity

type (
	// ChannelGroupsUpdated 频道组信息被更新时触发
	ChannelGroupsUpdated struct {
	}

	// ChannelsUpdated 频道信息被创建，更新，重新排序时触发
	ChannelsUpdated struct {
	}

	// ChannelDeleted 频道被删除时触发
	ChannelDeleted struct {
		Id int64 `json:"id"`
	}

	// ChannelMemberAdded 频道添加成员时触发
	ChannelMemberAdded struct {
		ChannelId int64 `json:"channelId"`
		UserId    int64 `json:"userId"`
	}

	// ChannelMemberRemoved 频道移除成员时触发
	ChannelMemberRemoved struct {
		ChannelId int64 `json:"channelId"`
		UserId    int64 `json:"userId"`
	}
)
