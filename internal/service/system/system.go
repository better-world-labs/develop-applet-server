package system

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"time"
)

//go:gone
func NewSystemService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Flag
	OssHost           string        `gone:"config,aliyun.oss.host"`
	OssAccessId       string        `gone:"config,aliyun.oss.access-key.id"`
	OssAccessSecret   string        `gone:"config,aliyun.oss.access-key.secret"`
	OssUseDir         string        `gone:"config,aliyun.oss.dir"`
	OssTokenExpiresIn time.Duration `gone:"config,aliyun.oss.expiresIn"`
	CronEmoticonStat  string        `gone:"config,cron.emoticon.use.stat"`

	P             iPersistence `gone:"*"`
	logrus.Logger `gone:"gone-logger"`
}
