package message

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type persistence struct {
	gone.Flag
	xorm.Engine   `gone:"gone-xorm"`
	logrus.Logger `gone:"gone-logger"`
}

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

func (p *persistence) listRangedMessageIdsByUserId(messageIds []int64, userId int64) ([]int64, error) {
	var res []int64

	return res, p.Table("message_like").Select("message_id id").
		In("message_id", messageIds).
		Where("user_id = ?", userId).Find(&res)
}

func (p *persistence) countMessageLikes(messageIds []int64) ([]*entity.MessageLikeCount, error) {
	var res []*entity.MessageLikeCount

	return res, p.Table("message_like").Select("message_id id, count(user_id) `like`").
		In("message_id", messageIds).Where("is_like  = 1").
		GroupBy("message_id").
		Find(&res)
}

func (p *persistence) createOrUpdateLike(like *entity.MessageLike) error {
	_, err := p.Exec("insert message_like (message_id, user_id, is_like) values (?, ?, ?) "+
		"on duplicate key update is_like = if(now() > created_at, ?, is_like)", like.MessageId, like.UserId, like.IsLike, like.IsLike)
	return err
}

func (p *persistence) create(m *message.Message) error {
	record, err := messageToMessageRecord(m)
	if err != nil {
		return nil
	}

	return p.Transaction(func(session xorm.Interface) error {
		exists, err := session.Table("message_record").Where("send_id = ?", record.SendId).Exist()
		if err != nil {
			return err
		}

		if !exists {
			_, err = p.Table("message_record").Insert(record)
		}

		m.Id = record.Id
		return err
	})
}

func (p *persistence) checkExists(sendId string) (bool, error) {
	return p.Table("message_record").Where("send_id = ?", sendId).Exist()
}

func (p *persistence) Get(sendId string) (*message.Message, error) {
	var record entity.MessageRecord

	_, err := p.Table("message_record").Where("send_id = ?", sendId).Get(&record)
	if err != nil {
		return nil, err
	}

	return recordToMessage(&record)
}

func (p *persistence) listHistoryMessageBefore(channelId, fromId int64, size int) ([]*message.Message, error) {
	var res []*entity.MessageRecord

	//表情回复类消息，不单独展示，所以过滤掉
	session := p.Table("message_record").
		Where("channel_id = ? and i_msg_type != ?", channelId, message.ContentTypeEmoticonReply)
	if fromId > 0 {
		session.And("id < ?", fromId)
	}

	err := session.Desc("id").
		Limit(size, 0).
		Find(&res)
	if err != nil {
		return nil, err
	}

	return recordsToMessages(res), err
}

func (p *persistence) listHistoryMessageAfter(channelId, fromId int64, size int) ([]*message.Message, error) {
	var res []*entity.MessageRecord

	//表情回复类消息，不单独展示，所以过滤掉
	session := p.Table("message_record").
		Where("channel_id = ? and i_msg_type != ?", channelId, message.ContentTypeEmoticonReply)
	if fromId > 0 {
		session.And("id > ?", fromId)
	}

	err := session.Asc("id").
		Limit(size, 0).
		Find(&res)
	if err != nil {
		return nil, err
	}

	return recordsToMessages(res), err
}

func (p *persistence) GetLastMessageByChannelId(channelId int64) (*message.Message, error) {
	var record entity.MessageRecord
	exists, err := p.Table("message_record").Where("channel_id = ? ", channelId).OrderBy("id desc").Get(&record)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return recordToMessage(&record)
}

func (p *persistence) getById(id int64) (*message.Message, bool, error) {
	var record entity.MessageRecord
	exists, err := p.Table("message_record").Where("id = ? ", id).Get(&record)
	m, err := recordToMessage(&record)
	if err != nil {
		return nil, false, err
	}

	return m, exists, err
}

func (p *persistence) listByIds(ids []int64) ([]*message.Message, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var records []*entity.MessageRecord

	err := p.Table("message_record").In("id", ids).Find(&records)
	if err != nil {
		return nil, err
	}

	return recordsToMessages(records), err
}

func (p *persistence) listByIdsMap(ids []int64) (map[int64]*message.Message, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	records := make(map[int64]*entity.MessageRecord, 0)
	messages := make(map[int64]*message.Message, len(records))
	err := p.Table("message_record").In("id", ids).Find(&records)
	if err != nil {
		return nil, err
	}

	for k, v := range records {
		if m, err := recordToMessage(v); err == nil {
			messages[k] = m
		}
	}

	return messages, nil
}

func (p *persistence) GetRecordsSummary(fromId, channelId int64) (int64, error) {
	count, err := p.Table("message_record").Where("channel_id = ? and id > ?", channelId, fromId).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *persistence) GetMessageCntByUserId(userId int64) (int64, error) {
	return p.Table("message_record").Where("user_id = ?", userId).Count()
}

func (p *persistence) insertReplyForMessage(msgId int64, emoticonId int64, userId int64) {
	_, err := p.Exec(`
		update message_record
		set content = if( JSON_LENGTH(content, '$.reply') > 0,
				JSON_ARRAY_INSERT(content, '$.reply[0]', JSON_OBJECT('userId', ?,'emoticonId', ?)),
				JSON_SET(content,'$.reply', JSON_ARRAY(JSON_OBJECT('userId', ?,'emoticonId', ?)))
			)
		where
		    id = ?
    `, userId, emoticonId, userId, emoticonId, msgId)

	if err != nil {
		p.Errorf("insertReplyForMessage(%d, %d, %d) err:%v", msgId, emoticonId, userId, err)
	}
}
