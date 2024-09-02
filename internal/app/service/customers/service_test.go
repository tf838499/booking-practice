package customers

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"ticket-booking/internal/app/service/customers/automock"
)

type serviceMock struct {
	CustomerRedis *automock.MockCustomerRedis
	OrderRepo     *automock.MockOrderRepository
}

func buildServiceMock(ctrl *gomock.Controller) serviceMock {
	return serviceMock{
		CustomerRedis: automock.NewMockCustomerRedis(ctrl),
		OrderRepo:     automock.NewMockOrderRepository(ctrl),
	}
}
func buildService(mock serviceMock) *CustomerService {
	param := CustomerServiceParam{
		CustomerRedis: mock.CustomerRedis,
		OrderRepo:     mock.OrderRepo,
	}
	return NewCustomerService(context.Background(), param)
}

// nolint
func TestMain(m *testing.M) {
	// To avoid getting an empty object slice
	_ = faker.SetRandomMapAndSliceMinSize(2)

	// To avoid getting a zero random number
	_ = faker.SetRandomNumberBoundaries(1, 100)

	m.Run()
}
