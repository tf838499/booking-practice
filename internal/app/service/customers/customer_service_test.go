package customers

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
	"ticket-booking/internal/domain/common"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCustomerService_CreateOrder(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		param   PostgresDB.CreateOrderAndSetSeatIsSoldParams
		tickets []PostgresDB.Ticket
	}
	var args Args
	_ = faker.FakeData(&args)
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		SetupService func(t *testing.T) *CustomerService
		wantErr      bool
	}{
		{
			name: "ticket create success",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().CreateOrderAndSetSeatIsSoldTx(gomock.Any(), args.param).Return(args.tickets, nil)
				service := buildService(mock)
				return service
			},
			wantErr: false,
		},
		{
			name: "ticket not found  creat item",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().CreateOrderAndSetSeatIsSoldTx(gomock.Any(), args.param).Return(args.tickets, sql.ErrNoRows)
				service := buildService(mock)
				return service
			},
			wantErr: true,
		},
		{
			name: "ticket create fail",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().CreateOrderAndSetSeatIsSoldTx(gomock.Any(), args.param).Return(args.tickets, errors.New(common.ErrorCodeInternalProcess.Name))
				service := buildService(mock)
				return service
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		c := tt
		t.Run(tt.name, func(t *testing.T) {
			service := c.SetupService(t)
			param := OrederParam{}
			_, err := service.CreateOrder(context.Background(), param)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCustomerService_GetSeatsList(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		param PostgresDB.SelectSeatsListByEventAndSectionParams
		Seat  []PostgresDB.Seat
	}
	var args Args
	_ = faker.FakeData(&args)
	args.Seat = []PostgresDB.Seat{
		{
			ID:         1,
			Section:    "A",
			SeatNumber: 1,
			SeatStatus: true,
			Price:      1000.01,
			EventID:    1,
		},
		{
			ID:         2,
			Section:    "A",
			SeatNumber: 2,
			SeatStatus: true,
			Price:      1000.01,
			EventID:    1,
		},
	}
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		SetupService func(t *testing.T) *CustomerService
		wantErr      bool
	}{
		{
			name: "seat does exist",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().SelectSeatsListAndUpdateSeatStatusTx(gomock.Any(), args.param).Return(args.Seat, nil)
				service := buildService(mock)
				return service
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c := tt
		t.Run(tt.name, func(t *testing.T) {
			service := c.SetupService(t)
			param := SeatsParam{}
			_, err := service.GetSeatsList(context.Background(), param)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCustomerService_UpdateTicketStatus(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		param        PostgresDB.UpdateTicketStatusParams
		TicketStatus []PostgresDB.UpdateTicketStatusRow
	}
	var args Args
	_ = faker.FakeData(&args)
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	args.param = PostgresDB.UpdateTicketStatusParams{
		OrderTradeNo: "test_1",
		TicketStatus: 1,
	}
	tests := []struct {
		name         string
		SetupService func(t *testing.T) *CustomerService
		wantErr      bool
	}{
		{
			name: "seat does exist",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().UpdateTicketStatus(gomock.Any(), args.param).Return(args.TicketStatus, nil)
				service := buildService(mock)
				return service
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c := tt
		t.Run(tt.name, func(t *testing.T) {
			service := c.SetupService(t)
			param := UpdateTicketParm{
				OrderTradeNoParam: "test_1",
				TicketStatusParam: 1,
			}
			_, err := service.UpdateTicketStatusIsPay(context.Background(), param)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCustomerService_UpdateSeatsAvailableByOrderTradeNo(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		param PostgresDB.UpdateSeatStatusBySeatIDsParams
	}
	var args Args
	_ = faker.FakeData(&args)
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	args.param.SeatStatus = false
	args.param.Column2 = []int32{3, 2}
	tests := []struct {
		name         string
		SetupService func(t *testing.T) *CustomerService
		wantErr      bool
	}{
		{
			name: "update seat success",
			SetupService: func(t *testing.T) *CustomerService {
				mock := buildServiceMock(ctrl)
				mock.OrderRepo.EXPECT().UpdateSeatStatusBySeatIDs(gomock.Any(), args.param).Return(nil)
				service := buildService(mock)
				return service
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		c := tt
		t.Run(tt.name, func(t *testing.T) {
			service := c.SetupService(t)
			param := UpdateSeatsAvailableParm{
				SeatID: []int32{3, 2},
			}
			err := service.UpdateSeatsAvailableBySeatID(context.Background(), param)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
