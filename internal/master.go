package internal

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gone-io/emitter"
	rocket "github.com/gone-io/emitter/adapter/rocket_aliyun_v4"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"os"
)

//go:generate gone priest -s . -p $GOPACKAGE -f Priest -o priest.go
func MasterPriest(cemetery gone.Cemetery) error {
	env := os.Getenv("ENV")
	if env == "local" || env == "" {
		_ = emitter.LocalMQPriest(cemetery)
	} else {
		_ = emitter.Priest(cemetery)
		_ = rocket.Priest(cemetery)
	}

	_ = goner.RedisPriest(cemetery)
	_ = goner.XormPriest(cemetery)
	_ = goner.SchedulePriest(cemetery)
	_ = goner.GinPriest(cemetery)

	_ = Priest(cemetery)
	return nil
}
