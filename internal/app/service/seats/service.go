package seats

import (
	"context"

	"github.com/rs/zerolog"
)

type SeatService struct {
	seatRepo  SeatRepository
	seatRedis SeatRedis
	dataRepo  DataRepository
}

type SeatServiceParam struct {
	SeatRepo  SeatRepository
	SeatRedis SeatRedis
	DataRepo  DataRepository
}

func NewSeatService(_ context.Context, param SeatServiceParam) *SeatService {
	return &SeatService{
		seatRepo:  param.SeatRepo,
		seatRedis: param.SeatRedis,
		dataRepo:  param.DataRepo,
	}
}
func (c *SeatService) logger(ctx context.Context) *zerolog.Logger {
	l := zerolog.Ctx(ctx).With().Str("component", "seat-service").Logger()
	return &l
}
