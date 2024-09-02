package rabbitMQ

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (q *RabbitMQ) SendOrderNotPay(ctx context.Context, orderTradeNo string) error {
	// rotuer := q.OrderQueue.Name
	err := q.Client.PublishWithContext(ctx,
		directExhangeName, // exchange
		q.OrderQueue.Name, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(orderTradeNo),
			Expiration:  "10000", // 設置消息過期時間（毫秒）,
			Headers: amqp.Table{
				"orderTradeNo": orderTradeNo,
			},
		})

	// fmt.Println("success", orderTradeNo)
	return err
}

// func (q *RabbitMQ) ReceiverOrder(ctx context.Context, orderTradeNo map[string]int) error {
func (q *RabbitMQ) ReceiverOrder(ctx context.Context, orderTradeNo chan map[string]int, PaidOrder chan []string) {

	msgs, err := q.Client.Consume(
		queueName, // 队列名字
		"",        // 消费者标记
		false,     // 自动应答
		true,      // 非排他
		false,     // 不本地
		false,     // 没有等待
		nil,       // 参数
	)
	if err != nil {
		fmt.Println(err)
	}

	AlreadyPaid := []string{}
	// OrderTradeNoMap := map[string]int{}

	for {
		select {
		case OrderTradeNoMap := <-orderTradeNo:

			first := ""
			msgProcessing := true
			for msgProcessing {
				select {
				case d, ok := <-msgs:
					if !ok {
						// 通道已關閉
						msgProcessing = false
						break
					}

					// log.Printf("Received a message: %s", d.Body)
					orderTradeNo := string(d.Body)
					_, orderExists := OrderTradeNoMap[orderTradeNo]

					if first == "" {
						first = orderTradeNo
					} else if first == orderTradeNo {
						d.Nack(false, true)
						msgProcessing = false
						break
					}

					if orderExists {
						d.Ack(true)
						AlreadyPaid = append(AlreadyPaid, orderTradeNo)
					} else {
						d.Nack(false, true)
					}
				default:
					// 沒有消息可供處理，退出內部循環
					msgProcessing = false
				}
			}

			// 將 AlreadyPaid 发送到 PaidOrder
			if len(AlreadyPaid) > 0 {
				PaidOrder <- AlreadyPaid
				AlreadyPaid = []string{}
			}

		case <-ctx.Done():
			return
		}
	}
}
func (q *RabbitMQ) DlxOrder(ctx context.Context, ch chan string) error {

	msgs, err := q.Client.Consume(
		dlxQueueName, // 隊列名字
		"",           // 消費者標記
		true,         // 自動應答
		false,        // 非排他
		false,        // 不本地
		false,        // 沒有等待
		nil,          // 參數
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			ch <- string(d.Body)
		}
	}()
	return nil
}
