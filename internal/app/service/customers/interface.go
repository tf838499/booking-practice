package customers

import (
	"context"
	RedisCache "ticket-booking/internal/adapter/redisClient"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
)

//go:generate mockgen -destination automock/order_repository.go -package=automock . OrderRepository
type OrderRepository interface {
	CreateOrderAndSetSeatIsSoldTx(ctx context.Context, arg PostgresDB.CreateOrderAndSetSeatIsSoldParams) ([]PostgresDB.Ticket, error)
	SelectSeatsListAndUpdateSeatStatusTx(ctx context.Context, arg PostgresDB.SelectSeatsListByEventAndSectionParams) ([]PostgresDB.Seat, error)
	UpdateTicketStatus(ctx context.Context, arg PostgresDB.UpdateTicketStatusParams) ([]PostgresDB.UpdateTicketStatusRow, error)
	UpdateSeatStatusBySeatIDs(ctx context.Context, arg PostgresDB.UpdateSeatStatusBySeatIDsParams) error
	UpdateSeatStatusByOrderTradeNo(context.Context, PostgresDB.UpdateSeatStatusByOrderTradeNoParams) ([]PostgresDB.Seat, error)
	DeletTicketAndUpdateSeatIsAvailableStatusTx(ctx context.Context, orderTradeNo string, seatStaus bool) ([]PostgresDB.Seat, error)
}

//go:generate mockgen -destination automock/customer_redis.go -package=automock . CustomerRedis
type CustomerRedis interface {
	SetSeatKey(ctx context.Context, arg RedisCache.SeatKeyParams) (bool, error)
	DelSeatKey(ctx context.Context, arg RedisCache.SeatKeyParams) error
	DelSeatValue(ctx context.Context, arg []PostgresDB.Seat) error
	SetSeatAndIncrSeatNumber(ctx context.Context, sections []string, seat []PostgresDB.Seat) error
	SetOrderTradeNo(ctx context.Context, orderTradeNo string) error
	GetOrderTradeNoExist(ctx context.Context, orderTradeNo string) (bool, error)
	GetOrderTradeNoList(ctx context.Context) ([]string, error)
	DelOrderTradeNo(ctx context.Context, orderTradeNo []string) error
}
type OrderMQ interface {
	SendOrderNotPay(ctx context.Context, orderTradeNo string) error
	ReceiverOrder(ctx context.Context, orderTradeNo chan map[string]int, PaidOrder chan []string)
	DlxOrder(ctx context.Context, ch chan string) error
}
