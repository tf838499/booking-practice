package RedisCache

import (
	"context"
	"testing"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"

	"github.com/stretchr/testify/require"
)

func TestRedisReposity_SetGood(t *testing.T) {
	db := getTestRedis()
	repo := initRepository(t, db)
	type Args struct {
		GoodFakes SeatKeyParams
	}

	var args Args
	// _ = faker.FakeData(&args)
	args.GoodFakes = SeatKeyParams{
		Section:    "Section 1",
		SeatNumber: []int32{1, 2},
	}

	_, err := repo.SetSeatKey(context.Background(), args.GoodFakes)

	require.NoError(t, err)

}

func TestRedisReposity_SetAndGetSeatValue(t *testing.T) {
	db := getTestRedis()
	repo := initRepository(t, db)
	type Args struct {
		GoodFakes []PostgresDB.Seat
	}

	var args Args
	// _ = faker.FakeData(&args)
	args.GoodFakes = []PostgresDB.Seat{
		{
			ID:         1,
			Section:    "Section 1",
			SeatNumber: 1,
			SeatStatus: true,
			Price:      100,
			EventID:    1,
		},
		{
			ID:         1,
			Section:    "Section 1",
			SeatNumber: 2,
			SeatStatus: true,
			Price:      100,
			EventID:    1,
		},
		{
			ID:         1,
			Section:    "Section 1",
			SeatNumber: 3,
			SeatStatus: true,
			Price:      100,
			EventID:    1,
		},
	}
	// []PostgresDB.Seat
	err := repo.SetSeatValueAndSectionAndIncrSeatNumber(context.Background(), []string{}, []float64{}, args.GoodFakes)
	args.GoodFakes = []PostgresDB.Seat{
		{
			ID:         1,
			Section:    "Section 1",
			SeatNumber: 3,
			SeatStatus: true,
			Price:      100,
			EventID:    1,
		},
	}
	repo.DelSeatValue(context.Background(), args.GoodFakes)
	// para:=[]GetSeatParams{{Section: "Section 1"}}
	repo.GetSeatValue(context.Background(), GetSeatParams{EventId: 1, Section: []string{"Section 1"}})
	require.NoError(t, err)

}
