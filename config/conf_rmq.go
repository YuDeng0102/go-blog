package config

type RabbitMQ struct {
	Dial string `json:"dial" yaml:"dial"` // RabbitMQ 连接地址
	// 例如 "amqp://guest:guest@localhost:5672/"
	// 其中 "guest:guest" 是用户名和密码，"localhost" 是 RabbitMQ 服务器地址，"5672" 是端口号
}
