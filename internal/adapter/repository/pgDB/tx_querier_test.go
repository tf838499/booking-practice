package PostgresDB

import (
	"context"
	"fmt"
	"testing"

	// "github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"ticket-booking/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_CreateOrderAndSetSeatIsSoldTx(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
	)

	// Args
	fakeData := CreateOrderAndSetSeatIsSoldParams{
		ConsumerID:   1,
		EventID:      1,
		SeatID:       []int32{1, 2, 3},
		TicketStatus: 3,
		OrderTradeNo: "ORDER_TRADE_NO_1",
	}

	data, err := repo.CreateOrderAndSetSeatIsSoldTx(context.Background(), fakeData)
	require.NoError(t, err)
	assert.Len(t, data, 3)
	seat1, err := repo.SelectSeat(context.Background(), 1)
	require.NoError(t, err)
	seat2, err := repo.SelectSeat(context.Background(), 2)
	require.NoError(t, err)
	seat3, err := repo.SelectSeat(context.Background(), 3)
	require.NoError(t, err)

	assert.Equal(t, false, seat1.SeatStatus)
	assert.Equal(t, false, seat2.SeatStatus)
	assert.Equal(t, false, seat3.SeatStatus)
	fmt.Println(data)
	// assertOrder(t, &fakeData, &data)
}

func TestPostgresRepository_SelectSeatsListAndUpdateSeatStatusTx(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
	)

	// Args

	data, err := repo.SelectSeatsListAndUpdateSeatStatusTx(context.Background(), SelectSeatsListByEventAndSectionParams{
		Section: "A",
		EventID: 1,
		Column3: []int32{1, 2, 3},
	})
	require.NoError(t, err)
	assert.Len(t, data, 3)
}
