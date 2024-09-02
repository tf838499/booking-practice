package seats

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"ticket-booking/internal/app/service/seats/automock"
)

type serviceMock struct {
	SeatRedis *automock.MockSeatRedis
	SeatRepo  *automock.MockSeatRepository
}

func buildServiceMock(ctrl *gomock.Controller) serviceMock {
	return serviceMock{
		SeatRedis: automock.NewMockSeatRedis(ctrl),
		SeatRepo:  automock.NewMockSeatRepository(ctrl),
	}
}
func buildService(mock serviceMock) *SeatService {
	param := SeatServiceParam{
		SeatRedis: mock.SeatRedis,
		SeatRepo:  mock.SeatRepo,
	}
	return NewSeatService(context.Background(), param)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid getting an empty object slice
	_ = faker.SetRandomMapAndSliceMinSize(2)

	// To avoid getting a zero random number
	_ = faker.SetRandomNumberBoundaries(1, 100)

	m.Run()
}
