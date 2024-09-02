package RedisCache

import (
	"context"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
)

type RedisQuerier interface {
	SetSeatKey(ctx context.Context, arg SeatKeyParams) (bool, error)
	SetSeatValueAndSectionAndIncrSeatNumber(ctx context.Context, sections []string, prices []float64, seat []PostgresDB.Seat) error
	SetSeatAndIncrSeatNumber(ctx context.Context, sections []string, seat []PostgresDB.Seat) error
	GetSectionsAndPricesByEventID(ctx context.Context, eventid int32) ([]PostgresDB.SelectSectionsAndPricesByEventIDRow, error)
	GetSectionsPricesBySectionsName(ctx context.Context, eventid int32, section string) (float64, error)
	GetSeatValue(ctx context.Context, arg GetSeatParams) ([]PostgresDB.Seat, error)
	DelSeatKey(ctx context.Context, arg SeatKeyParams) error
	GetSectionSeatCount(ctx context.Context, arg GetSectionSeatCount) (map[string]int32, error)
	IncrEventSeatNumberByEventID(ctx context.Context, eventid int32)
	DelSeatValue(ctx context.Context, arg []PostgresDB.Seat) error
	FlushAll(ctx context.Context, arg DelSeatParams)
	SetOrderTradeNo(ctx context.Context, orderTradeNo string) error
	GetOrderTradeNoExist(ctx context.Context, orderTradeNo string) (bool, error)
	GetOrderTradeNoList(ctx context.Context) ([]string, error)
	DelOrderTradeNo(ctx context.Context, orderTradeNo []string) error
}

// const (
// 	prefixCustomer = "Customer:Cart:"
// 	prefixPrice    = "Good:Price:"
// 	prefixUser     = "Oauth:User:"
// 	// store    = "store:"
// 	// good     = "good:"
// )

var _ RedisQuerier = (*RedisRepository)(nil)
