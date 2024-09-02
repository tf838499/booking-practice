package customers

import (
	"context"

	"github.com/rs/zerolog"
)

type CustomerService struct {
	orderRepo     OrderRepository
	customerRedis CustomerRedis
	orderMQ       OrderMQ
}

type CustomerServiceParam struct {
	OrderRepo     OrderRepository
	CustomerRedis CustomerRedis
	OrderMQ       OrderMQ
}

func NewCustomerService(_ context.Context, param CustomerServiceParam) *CustomerService {
	return &CustomerService{
		orderRepo:     param.OrderRepo,
		customerRedis: param.CustomerRedis,
		orderMQ:       param.OrderMQ,
	}
}
func (c *CustomerService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "customer-service").Logger()
	return &l
}
