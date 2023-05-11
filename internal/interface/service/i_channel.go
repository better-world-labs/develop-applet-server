package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type (
	IChannel interface {
		ListJoinedChannelIds(userId int64) ([]int64, error)

		ListChannels(planetId int64) ([]*entity.Channel, error)

		ListChannelsByIds(ids []int64) ([]*entity.Channel, error)

		ListChannelGroups(planetId int64) ([]*entity.ChannelGroup, error)

		ListChannelMembersSimple(channelId int64) ([]*entity.ChannelMember, error)

		ListChannelMembers(channelId int64) (*domain.ChannelMemberRes, error)

		ListMembersOffset(channelId int64, userIds []int64) (map[int64]int64, error)

		CreateNormalChannel(name string, icon string, groupId int64, planetId int64, createdBy int64, mute bool) (int64, error)

		CreatePrivateChannel(name string, icon string, groupId int64, planetId int64, createdBy int64, mute bool, expiresIn time.Duration) (int64, error)

		UpdateChannelName(channelGroupId int64, name string) error

		UpdateNotice(userId, channelId int64, notice string) error

		CreateChannelGroup(creatorId, planetId int64, name, icon string) (*domain.GroupIdRes, error)

		UpdateChannelGroupName(channelGroupId int64, name string) error

		DeleteChannel(id int64) error

		DeleteChannelGroup(id int64) error

		IsChannelMember(channelId, userId int64) (bool, error)

		AddChannelMember(channelId, userId int64) error

		RemoveChannelMember(channelId, userId int64) error

		GetChannelMember(channelId, userId int64) (*entity.ChannelMember, error)

		ApplyPrivateChannel(userId, channelId int64, reason string) error

		UpdateChannelSort(groupId int64, sort []int64) error

		UpdateChannelGroupSort(planetId int64, sort []int64) error

		GetChannelGroupByGroupId(groupId int64) (*entity.ChannelGroup, error)

		GetChannelByChannelId(channelId int64) (*entity.Channel, bool, error)

		IsChannelValid(channelId int64) (bool, error)

		UpdateLastReadMessage(userId, channelId int64) (int64, error)

		UpdateMemberState(userId, channelId int64, state entity.ChannelMemberState) error

		GetLastReadMsgId(userId, channelId int64) (*entity.ChannelMember, error)

		CheckSendUnreadMsgEvent(lastReadId, channelId int64) (bool, error)

		GetUnreadMsgNum(lastReadId, userId, channelId int64) (int64, error)

		GetChannelIdsByUserId(userId int64) ([]int64, error)

		GetAllUnReadCountByUser(userId int64, planetId int64) ([]*domain.UserChannelUnreadNum, error)

		IncreaseUnreadNum(channelId, incr, maxMsgId int64, exceptUserIds []int64) error
	}
)
