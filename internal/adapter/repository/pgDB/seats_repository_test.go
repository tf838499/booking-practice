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

func TestPostgresRepository_GetGoodByID(t *testing.T) {
	db := getTestPostgresDB()
	repo := initRepository(t, db,
		testdata.Path(testdata.TestDataSeat),
	)

	event_id := 1
	data, err := repo.SelectSeatsList(context.Background(), int32(event_id))
	fmt.Println(data)
	// assertGood(t, args.data, data)
	require.NoError(t, err)
	assert.Len(t, data, 8)
}
