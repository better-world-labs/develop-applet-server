package channel

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	IChannelPersistence interface {
		create(channel *entity.Channel) error
		listByPlanetId(planetId int64) ([]*entity.Channel, error)
		delete(id int64) error
		updateSort(groupId, channelId, sort int64) error
		get(channelId int64) (*entity.Channel, bool, error)
		updateName(channelId int64, name string) error
		updateNotice(channelId int64, notice string) error
		listByIds(ids []int64) ([]*entity.Channel, error)
	}

	IChannelMemberPersistence interface {
		list(channelId int64) ([]*entity.ChannelMember, error)
		listIdsByChannelId(channelId int64) ([]int64, error)
		listChannelIdsByUserId(userId int64) ([]int64, error)
		listMembersOffset(channelId int64, userIds []int64) (map[int64]int64, error)
		isExists(channelId, userId int64) (bool, error)
		get(channelId, userId int64) (*entity.ChannelMember, error)
		create(channel *entity.ChannelMember) error
		delete(id int64) error
		updateState(id int64, state entity.ChannelMemberState) error
		UpdateLastReadMsgId(userId, channelId, messageId int64) error
		BatchUpdateLastReadMsgId(userIds []int64, channelId, messageRecordId int64) error
		GetChannelIdsByUserId(userId int64) ([]int64, error)
	}

	IChannelGroupPersistence interface {
		listByPlanetId(planetId int64) ([]*entity.ChannelGroup, error)
		create(group *entity.ChannelGroup) error
		delete(id int64) error
		updateSort(groupId, channelId, sort int64) error
		get(groupId int64) (*entity.ChannelGroup, error)
		updateName(groupId int64, name string) error
	}
)
