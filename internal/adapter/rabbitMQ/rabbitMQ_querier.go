package rabbitMQ

import "context"

type RabbiqmqQuerier interface {
	SendOrderNotPay(ctx context.Context, orderTradeNo string) error
	ReceiverOrder(ctx context.Context, orderTradeNo chan map[string]int, PaidOrder chan []string)
	DlxOrder(ctx context.Context, ch chan string) error
}

var _ RabbiqmqQuerier = (*RabbitMQ)(nil)
