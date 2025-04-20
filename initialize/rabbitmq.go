package initialize

import (
	"server/global"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func RabbitmqInit() *amqp.Connection {
	conn, err := amqp.Dial(global.Config.RabbitMQ.Dial)
	if err != nil {
		global.Log.Error("Failed to connect to RabbitMQ:", zap.Error(err))
		return nil
	}
	return conn
}
