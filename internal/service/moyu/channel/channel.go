package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/redis"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"strconv"
	"time"
)

const (
	DefaultUnreadHashKey = "h_unread_channel_%d"
)

//go:gone
func NewChannelService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Flag
	PChannel       IChannelPersistence       `gone:"*"`
	PChannelGroup  IChannelGroupPersistence  `gone:"*"`
	PChannelMember IChannelMemberPersistence `gone:"*"`

	Approval             service.IApproval      `gone:"*"`
	UserService          service.IUser          `gone:"*"`
	Planet               service.IPlanet        `gone:"*"`
	MessageRecordService service.IMessageRecord `gone:"*"`
	RedisService         service.IRedisService  `gone:"*"`

	emitter.Sender `gone:"gone-emitter"`
	Logger         logrus.Logger `gone:"gone-logger"`

	xorm.Engine `gone:"gone-xorm"`
}

func (s *svc) ListJoinedChannelIds(userId int64) ([]int64, error) {
	return s.PChannelMember.listChannelIdsByUserId(userId)
}

func (s *svc) ListChannelsByIds(ids []int64) ([]*entity.Channel, error) {
	return s.PChannel.listByIds(ids)
}

func (s *svc) ListChannels(planetId int64) ([]*entity.Channel, error) {
	return s.PChannel.listByPlanetId(planetId)
}

func (s *svc) ListChannelGroups(planetId int64) ([]*entity.ChannelGroup, error) {
	return s.PChannelGroup.listByPlanetId(planetId)
}

func (s *svc) ListMembersOffset(channelId int64, userIds []int64) (map[int64]int64, error) {
	return s.PChannelMember.listMembersOffset(channelId, userIds)
}

// ListChannelMembersSimple TODO 和 ListChannelMembers 合并复用
func (s *svc) ListChannelMembersSimple(channelId int64) ([]*entity.ChannelMember, error) {
	return s.PChannelMember.list(channelId)
}

func (s *svc) ListChannelMembers(channelId int64) (*domain.ChannelMemberRes, error) {
	userIds, err := s.PChannelMember.listIdsByChannelId(channelId)
	if err != nil {
		return nil, err
	}
	s.Logger.Infof("ListChannelMembers: channelId=%d len=%d", channelId, len(userIds))

	channel, exists, err := s.PChannel.get(channelId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, gin.NewParameterError("channel not found")
	}

	roles, err := s.Planet.ListPlanetRolesMap(channel.PlanetId, userIds)
	if err != nil {
		return nil, err
	}

	members := make([]*domain.ChannelMember, 0)
	users, err := s.UserService.GetUserInBatch(userIds)
	if err != nil {
		return nil, err
	}

	var online int64
	for _, aUser := range users {
		if aUser.Online {
			online++
		}

		members = append(members, &domain.ChannelMember{
			Id:       aUser.Id,
			Nickname: aUser.Nickname,
			Avatar:   string(aUser.Avatar),
			Online:   aUser.Online,
			Role:     roles[aUser.Id],
		})
	}

	return &domain.ChannelMemberRes{
		Total:          int64(len(members)),
		Online:         online,
		ChannelMembers: members,
	}, err
}

func (s *svc) CreateNormalChannel(name string, icon string, groupId int64, planetId int64, createdBy int64, mute bool) (int64, error) {
	return s.createChannel(&entity.Channel{
		Name:      name,
		Icon:      icon,
		Type:      entity.ChannelTypeNormal,
		GroupId:   groupId,
		PlanetId:  planetId,
		Mute:      mute,
		State:     1,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		Sort:      time.Now().Unix(), // use timestamp as current sort
	})
}

func (s *svc) CreatePrivateChannel(name string, icon string, groupId int64, planetId int64, createdBy int64, mute bool, expiresIn time.Duration) (int64, error) {
	expiresInTime := time.Now().Add(expiresIn)
	return s.createChannel(&entity.Channel{
		Name:      name,
		Icon:      icon,
		Type:      entity.ChannelTypePrivate,
		Mute:      mute,
		State:     1,
		GroupId:   groupId,
		PlanetId:  planetId,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		ExpiresAt: &expiresInTime,
		Sort:      time.Now().Unix(), // use timestamp as current sort
	})
}

func (s *svc) createChannel(channel *entity.Channel) (int64, error) {
	err := s.PChannel.create(channel)
	if err != nil {
		return 0, err
	}

	err = s.Send(&entity.ChannelsUpdated{})
	if err != nil {
		return 0, err
	}

	return channel.Id, err
}

func (s *svc) UpdateNotice(userId, channelId int64, notice string) error {
	err := s.PChannel.updateNotice(channelId, notice)
	if err != nil {
		return err
	}

	id := uuid.NewString()
	now := time.Now()
	b, err := json.Marshal(message.NewChannelNoticeContent(0, nil, notice))
	if err != nil {
		return err
	}

	if len(notice) == 0 {
		return nil
	}

	return s.Sender.Send(&wsevent.MsgSendEvent{
		ClientId: id,
		Id:       id,
		UserId:   userId,
		SeqId:    0,
		Msg: wsevent.Msg{
			CreatedAt: now,
			UserId:    userId,
			ChannelId: channelId,
			Content:   b,
		},
		CreatedAt: now,
	})
}

func (s *svc) UpdateChannelName(channelId int64, name string) error {
	err := s.PChannel.updateName(channelId, name)
	if err != nil {
		return err
	}

	return s.Send(&entity.ChannelsUpdated{})
}

func (s *svc) CreateChannelGroup(creatorId, planetId int64, name, icon string) (*domain.GroupIdRes, error) {
	group := entity.ChannelGroup{
		Name:      name,
		Icon:      icon,
		PlanetId:  planetId,
		CreatedBy: creatorId,
		CreatedAt: time.Now(),
	}
	err := s.PChannelGroup.create(&group)
	if err != nil {
		return nil, err
	}

	return &domain.GroupIdRes{GroupId: group.Id}, s.Send(&entity.ChannelGroupsUpdated{})
}

func (s *svc) UpdateChannelGroupName(channelGroupId int64, name string) error {
	err := s.PChannelGroup.updateName(channelGroupId, name)
	if err != nil {
		return err
	}

	return s.Send(&entity.ChannelsUpdated{})
}

func (s *svc) DeleteChannel(id int64) error {
	err := s.PChannel.delete(id)
	if err != nil {
		return err
	}

	return s.Send(&entity.ChannelDeleted{Id: id})
}

func (s *svc) DeleteChannelGroup(id int64) error {
	err := s.PChannelGroup.delete(id)
	if err != nil {
		return err
	}

	return s.Send(&entity.ChannelGroupsUpdated{})
}

func (s *svc) UpdateChannelGroupSort(planetId int64, sort []int64) error {
	type orderPair struct {
		Sort           int64
		ChannelGroupId int64
	}

	pairs := make([]orderPair, 0)
	for i, order := range sort {
		pairs = append(pairs, orderPair{
			Sort:           int64(i),
			ChannelGroupId: order,
		})
	}

	for _, pair := range pairs {
		err := s.PChannelGroup.updateSort(planetId, pair.ChannelGroupId, pair.Sort)
		if err != nil {
			return err
		}
	}

	return s.Send(&entity.ChannelGroupsUpdated{})
}

func (s *svc) UpdateChannelSort(groupId int64, sort []int64) error {
	type orderPair struct {
		Sort      int64
		ChannelId int64
	}

	pairs := make([]orderPair, 0)
	for i, order := range sort {
		pairs = append(pairs, orderPair{
			Sort:      int64(i),
			ChannelId: order,
		})
	}

	for _, pair := range pairs {
		err := s.PChannel.updateSort(groupId, pair.ChannelId, pair.Sort)
		if err != nil {
			return err
		}
	}

	return s.Send(&entity.ChannelsUpdated{})
}

func (s *svc) IsChannelMember(channelId, userId int64) (bool, error) {
	return s.PChannelMember.isExists(channelId, userId)
}

func (s *svc) AddChannelMember(channelId, userId int64) error {
	return s.joinChannel(userId, channelId)
}

func (s *svc) GetChannelGroupByGroupId(groupId int64) (*entity.ChannelGroup, error) {
	return s.PChannelGroup.get(groupId)
}

func (s *svc) IsChannelValid(channelId int64) (bool, error) {
	channel, exists, err := s.GetChannelByChannelId(channelId)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	return channel.Status() == entity.ChannelStatusValid, nil
}

func (s *svc) GetChannelByChannelId(channelId int64) (*entity.Channel, bool, error) {
	return s.PChannel.get(channelId)
}

func (s *svc) joinChannel(userId, channelId int64) error {
	exists, err := s.IsChannelMember(channelId, userId)
	if err != nil {
		s.Logger.Errorf("Is channel member check failed, err: %s", err.Error())
		return err
	}

	if exists {
		return nil
	}

	s.Logger.Infof("joinChannel: channelId=%d, userId=%d", channelId, userId)
	err = s.PChannelMember.create(&entity.ChannelMember{
		ChannelId: channelId,
		UserId:    userId,
		State:     entity.ChannelMemberStateJoined,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	get, err := s.PChannelMember.get(channelId, userId)
	if err != nil {
		return err
	}

	s.Logger.Infof("create member: %v", get)

	return s.Send(&entity.ChannelMemberAdded{ChannelId: channelId, UserId: userId})
}

func (s *svc) GetChannelMember(channelId, userId int64) (*entity.ChannelMember, error) {
	return s.PChannelMember.get(channelId, userId)
}

func (s *svc) RemoveChannelMember(channelId, userId int64) error {
	err := s.UpdateMemberState(userId, channelId, entity.ChannelMemberStateRemoved)
	if err != nil {
		return err
	}

	return s.Send(&entity.ChannelMemberRemoved{ChannelId: channelId, UserId: userId})
}

func (s *svc) checkAdmin(userId, channelId int64) (bool, error) {
	channel, exists, err := s.PChannel.get(channelId)
	if err != nil {
		return exists, err
	}

	if !exists {
		return exists, gin.NewParameterError("channel not found")
	}

	role, err := s.Planet.GetPlanetRoles(channel.PlanetId, userId)
	if err != nil {
		return false, err
	}

	return role >= entity.PlanetRoleAdmin, nil
}

func (s *svc) ApplyPrivateChannel(userId, channelId int64, reason string) error {
	isAdmin, err := s.checkAdmin(userId, channelId)
	if err != nil {
		return err
	}

	if isAdmin {
		return s.AddChannelMember(channelId, userId)
	}

	return s.Transaction(func(session xorm.Interface) error {
		member, err := s.PChannelMember.get(channelId, userId)
		if err != nil {
			return err
		}

		if member != nil {
			if member.State == entity.ChannelMemberStateJoined {
				return errors.New("already in this channel")
			}

			if member.State == entity.ChannelMemberStateApplying {
				return nil
			}
		}

		channel, exists, err := s.PChannel.get(channelId)
		if err != nil {
			return err
		}

		if !exists {
			return gin.NewParameterError("channel not found")
		}

		if channel.Status() == entity.ChannelStatusInvalid {
			return errors.New("invalid channel status")
		}

		approval, err := s.Approval.StartApprove(entity.ApprovalTypeChannelJoin, userId, reason, channelId)
		if err != nil {
			return err
		}

		return s.PChannelMember.create(&entity.ChannelMember{
			ChannelId: channelId,
			UserId:    userId,
			ApplyId:   approval.Id,
			State:     entity.ChannelMemberStateApplying,
			CreatedAt: time.Now(),
		})
	})
}

func (s *svc) UpdateMemberState(userId, channelId int64, state entity.ChannelMemberState) error {
	member, err := s.PChannelMember.get(channelId, userId)
	if err != nil {
		return nil
	}

	if member != nil {
		return s.PChannelMember.updateState(member.Id, state)
	}

	s.Logger.Warnf("UpdateMemberState: member not found. userId=%d, channelId=%d", userId, channelId)
	return nil
}

func (s *svc) UpdateLastReadMessage(userId, channelId int64) (int64, error) {
	messageRecord, err := s.MessageRecordService.GetLastMessageByChannelId(channelId)
	if err != nil {
		s.Logger.Errorf("get channel's last message record failed, err: %s", err.Error())
		return 0, err
	}

	if messageRecord == nil {
		s.Logger.Infof("channel[%d] has no message record found. ", channelId)
		return 0, nil
	}

	err = s.PChannelMember.UpdateLastReadMsgId(userId, channelId, messageRecord.Id)
	if err != nil {
		s.Logger.Errorf("update last message record failed, err: %s", err.Error())
		err = gone.NewInnerError(500, err.Error())
	}
	// update 0 to redis, use default planet id [1] at present
	key := fmt.Sprintf(DefaultUnreadHashKey, channelId)
	field := fmt.Sprintf("%d", userId)
	s.RedisService.HSet(key, field, "0")

	return messageRecord.Id, err
}

func (s *svc) GetLastReadMsgId(userId, channelId int64) (*entity.ChannelMember, error) {
	return s.PChannelMember.get(channelId, userId)
}

func (s *svc) CheckSendUnreadMsgEvent(lastReadId, channelId int64) (bool, error) {
	records, err := s.MessageRecordService.ListHistory(channelId, lastReadId, 11, false)
	if err != nil {
		s.Logger.Errorf("get history records failed, err: %s", err.Error())
		return false, err
	}

	// 只要有未读消息就应该发送事件
	if len(records) > 0 {
		return true, nil
	}

	return false, nil
}

func (s *svc) GetUnreadMsgNum(lastReadId, userId, channelId int64) (int64, error) {
	key := fmt.Sprintf(DefaultUnreadHashKey, channelId)
	field := fmt.Sprintf("%d", userId)

	var total int64
	value, err := s.RedisService.HGet(key, field)
	if err != nil {
		s.Logger.Infof("Get unread cnt failed from redis, trying to get from db. last-read %d user-id %d channel-id %d", lastReadId, userId, channelId)
		if err != redis.ErrNil {
			s.Logger.Errorf("Use Command [HGET] at Redis failed. err: %s", err.Error())
			return 0, err
		}

		total, err = s.MessageRecordService.GetRecordsNum(lastReadId, channelId)
		if err != nil {
			s.Logger.Errorf("get unread records total failed, err: %s", err.Error())
			return 0, err
		}

		// write back to redis
		s.RedisService.HSet(key, field, strconv.FormatInt(total, 10))
		value = []byte(strconv.FormatInt(total, 10))
	}

	total, err = strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		s.Logger.Errorf("redis value [%s] of key [%s], filed [%s]  is not a number. ", value, key, field)
		total = 0
	}

	s.Logger.Infof("%d unread msg found. user-id %d, channel-id %d", total, userId, channelId)
	return total, nil
}

func (s *svc) GetChannelIdsByUserId(userId int64) ([]int64, error) {
	return s.PChannelMember.GetChannelIdsByUserId(userId)
}

func (s *svc) GetAllUnReadCountByUser(userId int64, planetId int64) ([]*domain.UserChannelUnreadNum, error) {
	var unreadNums []*domain.UserChannelUnreadNum

	channelIds, err := s.PChannelMember.GetChannelIdsByUserId(userId)
	if err != nil {
		s.Logger.Errorf("get user's channel failed. err: %s", err.Error())
		return nil, err
	}

	for _, channelId := range channelIds {
		key := fmt.Sprintf(DefaultUnreadHashKey, channelId)
		field := fmt.Sprintf("%d", userId)
		reply, err := s.RedisService.HGet(key, field)
		if err != nil {
			if err != redis.ErrNil {
				s.Logger.Errorf("Use Command [HGET] at Redis failed. err: %s", err.Error())
				return nil, err
			}

			s.Logger.Infof("Hash not contains filed [%s] of key [%s], trying to get from db.", field, key)

			channelMember, err := s.PChannelMember.get(channelId, userId)
			if err != nil {
				s.Logger.Errorf("Get user [%d] at channel [%d] last read msg id failed. err: %s", userId, channelId, err.Error())
				return nil, err
			}

			if channelMember == nil {
				s.Logger.Infof("User [%d] not in channel [%d], skip it. ", userId, channelId)
				continue
			}

			total, err := s.MessageRecordService.GetRecordsNum(channelMember.LastReadMessageId, channelId)
			if err != nil {
				s.Logger.Errorf("Get user [%d] at channel [%d] unread msg num failed. err: %s", userId, channelId, err.Error())
				return nil, err
			}

			reply = []byte(strconv.FormatInt(total, 10))

			// write back to Redis
			s.RedisService.HSet(key, field, string(reply))
		}

		num, err := strconv.ParseInt(string(reply), 10, 64)
		if err != nil {
			s.Logger.Errorf("redis value [%s] of key [%s], filed [%s]  is not a number. ", reply, key, field)
			num = 0
		}

		unreadNums = append(unreadNums, &domain.UserChannelUnreadNum{
			ChannelId: channelId,
			UnreadNum: num,
		})
	}

	return unreadNums, nil
}

func (s *svc) IncreaseUnreadNum(channelId, incr, maxMsgId int64, exceptUserIds []int64) error {
	userIds, err := s.PChannelMember.listIdsByChannelId(channelId)

	linq.From(exceptUserIds).Distinct().ToSlice(&exceptUserIds)

	linq.From(userIds).WhereT(func(userId int64) bool {
		for _, ex := range exceptUserIds {
			if userId == ex {
				return false
			}
		}
		return true
	}).ToSlice(&userIds)

	//  1. update redis
	for _, userId := range userIds {
		key := fmt.Sprintf(DefaultUnreadHashKey, channelId)
		field := fmt.Sprintf("%d", userId)
		err = s.RedisService.HIncrBy(key, field, incr)
		if err != nil {
			s.Logger.Errorf("incr unread num failed. err: %s", err.Error())
			return err
		}
	}

	// 2. update last read msg id
	return s.PChannelMember.BatchUpdateLastReadMsgId(exceptUserIds, channelId, maxMsgId)
}
