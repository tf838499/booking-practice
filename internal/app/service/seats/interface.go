package seats

import (
	"context"
	RedisCache "ticket-booking/internal/adapter/redisClient"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"

	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination automock/seat_repository.go -package=automock . SeatRepository
type SeatRepository interface {
	SelectSeatsList(ctx context.Context, eventID int32) ([]PostgresDB.Seat, error)
	SelectSectionsAndPricesByEventID(ctx context.Context, eventID int32) ([]PostgresDB.SelectSectionsAndPricesByEventIDRow, error)
}
type DataRepository interface {
	CreateConsumer(ctx context.Context, arg PostgresDB.CreateConsumerParams) (PostgresDB.Consumer, error)
	CreateEvent(ctx context.Context, arg PostgresDB.CreateEventParams) (PostgresDB.Event, error)
	CreateSeat(ctx context.Context, arg PostgresDB.CreateSeatParams) (PostgresDB.Seat, error)
}

//go:generate mockgen -destination automock/seat_redis.go -package=automock . SeatRedis
type SeatRedis interface {
	SetSeatValueAndSectionAndIncrSeatNumber(ctx context.Context, sections []string, prices []float64, arg []PostgresDB.Seat) error
	GetSeatValue(ctx context.Context, arg RedisCache.GetSeatParams) ([]PostgresDB.Seat, error)
	GetSectionsAndPricesByEventID(ctx context.Context, eventid int32) ([]PostgresDB.SelectSectionsAndPricesByEventIDRow, error)
	GetSectionsPricesBySectionsName(ctx context.Context, eventid int32, section string) (float64, error)
	GetSectionSeatCount(ctx context.Context, arg RedisCache.GetSectionSeatCount) (map[string]int32, error)
	FlushAll(ctx context.Context, arg RedisCache.DelSeatParams)
}
