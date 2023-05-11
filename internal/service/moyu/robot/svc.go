package robot

import (
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Goner
	P   iConfigPersistence `gone:"*"`
	Env string             `gone:"config,server.env"`

	Message service.IMessageRecord `gone:"*"`
	User    service.IUser          `gone:"*"`
	Channel service.IChannel       `gone:"*"`
}

func (s svc) createRobotUser(openid string) (int64, error) {
	//TODO
	return 0, nil
}

func (s svc) CreateConfig(userIda, channelId int64, messageReceiveUrl string) error {
	//channel, exists, err := s.Channel.GetChannelByChannelId(channelId)
	//if err != nil {
	//	return err
	//}
	//
	//if !exists {
	//	return errors.New("channel not found")
	//}
	//
	//s.createRobotUser(,)
	//s.P.CreateConfig(&entity.RobotConfig{
	//	CreatedBy:            userId,
	//	ChannelId:         channelId,
	//	MessageReceiveUrl: messageReceiveUrl,
	//})
	//
	//TODO
	return nil
}

func (s svc) ListConfigsByTrigger(triggerType entity.TriggerType) ([]*entity.RobotConfig, error) {
	return s.P.ListByTriggerType(triggerType)
}

func (s svc) ListConfigs() ([]*entity.RobotConfig, error) {
	return s.P.List()
}

func (s svc) GetConfigByRobotId(appId string) (*entity.RobotConfig, bool, error) {
	return s.P.GetByRobotId(appId)
}

func (s svc) ListConfigsByUserIds(userIds []int64) ([]*entity.RobotConfig, error) {
	return s.P.ListByUserIds(userIds)
}

func (s svc) SendMessage(env, robotId string, channelId int64, content []byte) error {
	//TODO cache
	config, exists, err := s.P.GetByRobotId(robotId)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("config not found")
	}

	// config 指定了环境，则 env 必须是 config 指定的环境
	if len(config.Env) > 0 && config.Env != env {
		return nil
	}

	// config 没有指定环境，env 指定了环境，则只能是服务本身的环境
	if len(config.Env) == 0 && len(env) > 0 && env != s.Env {
		return nil
	}

	return s.Message.SendMessage(fmt.Sprintf("robot_%s", robotId), config.UserId, channelId, content)
}
