package rabbitMQ

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Client     *amqp.Channel
	OrderQueue *amqp.Queue
	DlxQueue   *amqp.Queue
}
type MqParam struct {
	NormalQueue    string
	MessageTTL     int
	DeadExchange   string
	DeadRoutingKey string
}

const (
	directExhangeName = "direct_exchange"
	dlxExhangeName    = "dlx_exchange"
	queueName         = "direct.queue1"
	dlxQueueName      = "dlxdirect.queue1"
)

func NewRabbitMQ(client *amqp.Channel) (*RabbitMQ, error) {
	var err error
	// 聲明死信交換機
	err = client.ExchangeDeclare(
		dlxExhangeName, // 交換機的名字
		"fanout",       // 交換機的類型
		true,           // 持久化
		false,          // 自動刪除
		false,          // 內部
		false,          // 沒有等待
		nil,            // 參數
	)
	if err != nil {
		return nil, err
	}

	// 聲明死信隊列
	dlxQueue, err := client.QueueDeclare(
		dlxQueueName, // 隊列的名字
		true,         // 持久化
		false,        // 自動刪除
		false,        // 排他性
		false,        // 沒有等待
		nil,          // 參數
	)
	if err != nil {
		return nil, err
	}

	// 綁定死信隊列到死信交換機
	err = client.QueueBind(
		dlxQueue.Name,  // 隊列名字
		"",             // 路由鍵
		dlxExhangeName, // 交換機名字
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	// 聲明主交換機
	err = client.ExchangeDeclare(
		directExhangeName, // 交換機的名字
		"direct",          // 交換機的類型
		true,              // 持久化
		false,             // 自動刪除
		false,             // 內部
		false,             // 沒有等待
		nil,               // 參數
	)
	if err != nil {
		return nil, err
	}
	// 聲明主隊列，並設置死信交換機參數
	args := amqp.Table{
		// "x-message-ttl":          int32(10000),
		"x-dead-letter-exchange": dlxExhangeName,
	}
	mainQueue, err := client.QueueDeclare(
		queueName, // 隊列的名字
		true,      // 持久化
		false,     // 自動刪除
		false,     // 排他性
		false,     // 沒有等待
		args,      // 參數
	)
	if err != nil {
		return nil, err
	}
	// 綁定主隊列到主交換機
	err = client.QueueBind(
		mainQueue.Name,    // 隊列名字
		mainQueue.Name,    // 路由鍵
		directExhangeName, // 交換機名字
		false,
		nil,
	)
	return &RabbitMQ{
		Client:     client,
		OrderQueue: &mainQueue,
		DlxQueue:   &dlxQueue,
	}, err
}
