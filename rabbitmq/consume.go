package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeMessages(conn *amqp.Connection, queueName string) (<-chan amqp.Delivery, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// 设置 QoS 控制并发
	err = ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("qos设置失败: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumerTag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	return msgs, err
}
