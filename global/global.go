package global

import (
	"server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"gorm.io/gorm"

	"go.uber.org/zap"
)

var (
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      redis.Client
	ESClient   *elasticsearch.TypedClient
	BlackCache local_cache.Cache
	RmqConn    *amqp091.Connection
)
