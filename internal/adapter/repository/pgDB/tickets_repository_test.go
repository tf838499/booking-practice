package PostgresDB

import (
	"context"
	"testing"
	"time"

	// "github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"ticket-booking/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertOrder(t *testing.T, expected *Ticket, actual *Ticket) {
	require.NotNil(t, actual)
	assert.Equal(t, expected.EventID, actual.EventID)
	assert.Equal(t, expected.ConsumerID, actual.ConsumerID)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.OrderTradeNo, actual.OrderTradeNo)
	assert.NotNil(t, actual.PurchaseDate)
	assert.Equal(t, expected.SeatID, actual.SeatID)
	assert.Equal(t, expected.TicketStatus, actual.TicketStatus)
}

func TestPostgresRepository_CreateOrder(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
	)

	// Args
	fakeData := Ticket{
		ID:           10001,
		ConsumerID:   1,
		EventID:      1,
		SeatID:       1,
		OrderTradeNo: "1",
		PurchaseDate: time.Now(),
		TicketStatus: 1,
	}
	data, err := repo.CreateOrder(context.Background(), CreateOrderParams{
		ConsumerID:   fakeData.ConsumerID,
		EventID:      fakeData.EventID,
		SeatID:       fakeData.SeatID,
		TicketStatus: fakeData.TicketStatus,
		OrderTradeNo: fakeData.OrderTradeNo,
	})
	require.NoError(t, err)
	assertOrder(t, &fakeData, &data)
}
func TestPostgresRepository_SelectAndUpdateSeatsListByEventAndSection(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
	)

	// Args

	data, err := repo.SelectSeatsListByEventAndSection(context.Background(), SelectSeatsListByEventAndSectionParams{
		Section: "A",
		EventID: 1,
		Column3: []int32{1, 2, 3},
	})
	require.NoError(t, err)
	assert.Len(t, data, 3)
}
func TestPostgresRepository_UpdateSeatStatus(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
	)

	// Args
	seatID := 1
	_, err := repo.UpdateSeatStatus(context.Background(), UpdateSeatStatusParams{
		ID:         int32(seatID),
		SeatStatus: false,
	})
	require.NoError(t, err)
	data, err := repo.SelectSeat(context.Background(), int32(seatID))
	require.NoError(t, err)
	assert.Equal(t, false, data.SeatStatus)
}

func TestPostgresRepository_UpdateSeatStatusBySeatIDs(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
		testdata.Path(testdata.TestDataTicket),
	)

	// Args
	seatID_1 := 1
	seatID_2 := 2
	err := repo.UpdateSeatStatusBySeatIDs(context.Background(), UpdateSeatStatusBySeatIDsParams{
		Column2:    []int32{1, 2},
		SeatStatus: false,
	})
	require.NoError(t, err)

	data_1, err := repo.SelectSeat(context.Background(), int32(seatID_1))
	require.NoError(t, err)
	assert.Equal(t, false, data_1.SeatStatus)

	data_2, err := repo.SelectSeat(context.Background(), int32(seatID_2))
	require.NoError(t, err)
	assert.Equal(t, false, data_2.SeatStatus)
}

func TestPostgresRepository_UpdateTicketStatus(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataConsumer),
		testdata.Path(testdata.TestDataEvent),
		testdata.Path(testdata.TestDataSeat),
		testdata.Path(testdata.TestDataTicket),
	)

	// Args

	ticketStatus := 3
	data, err := repo.UpdateTicketStatus(context.Background(), UpdateTicketStatusParams{
		OrderTradeNo: "ORDER_TRADE_NO_1",
		TicketStatus: int32(ticketStatus),
	})
	require.NoError(t, err)
	for i := range data {
		assert.Equal(t, ticketStatus, data[i].TicketStatus)
	}

}
