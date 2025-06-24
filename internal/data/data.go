package data

import (
	"layout/internal/conf"

	"github.com/go-redis/redis/extra/redisotel"
	redis "github.com/go-redis/redis/v8"
	"github.com/olaola-chat/psl-be-partystar-pkg/database"
	"github.com/olaola-chat/psl-be-partystar-pkg/metrics"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCouponRepo)

// Data .
type Data struct {
	XianshiDB *gorm.DB
	ConfigDB  *gorm.DB
	Rds       *redis.Client
}

// NewData .
func NewData(bc *conf.Bootstrap, logger log.Logger) (*Data, func(), error) {
	c, cs := bc.GetData(), bc.GetServer()
	d := &Data{}
	xianshi, err := database.NewMysqlDB(&database.MysqlConf{
		MaxOpenConns:    c.GetDatabase().GetXianshi().GetMaxOpenConns(),
		MaxIdleConns:    c.GetDatabase().GetXianshi().GetMaxIdleConns(),
		ConnMaxIdleTime: c.GetDatabase().GetXianshi().GetConnMaxIdleTime().AsDuration(),
		ConnMaxLifetime: c.GetDatabase().GetXianshi().GetConnMaxLifetime().AsDuration(),
		MasterDsn:       c.GetDatabase().GetXianshi().GetMaster(),
		SlaveDsn:        c.GetDatabase().GetXianshi().GetSlaves(),
	}, "partystar-server", cs.GetEnv())
	if err != nil {
		return nil, nil, err
	}
	d.XianshiDB = xianshi
	config, err := database.NewMysqlDB(&database.MysqlConf{
		MaxOpenConns:    c.GetDatabase().GetConfig().GetMaxOpenConns(),
		MaxIdleConns:    c.GetDatabase().GetConfig().GetMaxIdleConns(),
		ConnMaxIdleTime: c.GetDatabase().GetConfig().GetConnMaxIdleTime().AsDuration(),
		ConnMaxLifetime: c.GetDatabase().GetConfig().GetConnMaxLifetime().AsDuration(),
		MasterDsn:       c.GetDatabase().GetConfig().GetMaster(),
		SlaveDsn:        c.GetDatabase().GetConfig().GetSlaves(),
	}, "partystar-server", cs.GetEnv())
	if err != nil {
		return nil, nil, err
	}
	d.ConfigDB = config

	rds := redis.NewClient(&redis.Options{
		Network:      c.GetRedis().GetNetwork(),
		Addr:         c.GetRedis().GetAddr(),
		ReadTimeout:  c.GetRedis().GetReadTimeout().AsDuration(),
		WriteTimeout: c.GetRedis().GetWriteTimeout().AsDuration(),
		Password:     c.GetRedis().GetPassword(),
	})
	rds.AddHook(redisotel.TracingHook{})                                 // 链路追踪
	rds.AddHook(metrics.NewMetricsHook("partystar-server", cs.GetEnv())) // metrics

	d.Rds = rds

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		sd, _ := d.XianshiDB.DB()
		sd.Close()
		sd, _ = d.ConfigDB.DB()
		sd.Close()

		d.Rds.Close()
	}
	return d, cleanup, nil
}
