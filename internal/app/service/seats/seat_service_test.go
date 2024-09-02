package seats

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

func TestSeatService_ListSeats(t *testing.T) {
	t.Parallel()
	// Args
	type Args struct {
		param SeatListParam
		seats []PostgresDB.Seat
	}
	var args Args
	_ = faker.FakeData(&args)
	args.param.EventID = 1
	// Init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		SetupService func(t *testing.T) *SeatService
		wantErr      bool
	}{
		{
			name: "seat does exist",
			SetupService: func(t *testing.T) *SeatService {
				mock := buildServiceMock(ctrl)
				mock.SeatRepo.EXPECT().SelectSeatsList(gomock.Any(), args.param.EventID).Return(args.seats, nil)
				// mock.TraderRepo.EXPECT().GetTraderByEmail(gomock.Any(), args.Trader.Email).Return(nil, common.DomainError{})
				service := buildService(mock)
				return service
			},
			wantErr: false,
		},
		{
			name: "seat does not exist",
			SetupService: func(t *testing.T) *SeatService {
				mock := buildServiceMock(ctrl)
				seat := []PostgresDB.Seat{}
				mock.SeatRepo.EXPECT().SelectSeatsList(gomock.Any(), args.param.EventID).Return(seat, sql.ErrNoRows)
				service := buildService(mock)
				return service
			},
			wantErr: true,
		},
		{
			name: "failed to get good",
			SetupService: func(t *testing.T) *SeatService {
				mock := buildServiceMock(ctrl)
				seat := []PostgresDB.Seat{}
				mock.SeatRepo.EXPECT().SelectSeatsList(gomock.Any(), args.param.EventID).Return(seat, errors.New(common.ErrorCodeInternalProcess.Name))
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
			param := SeatListParam{
				EventID: args.param.EventID,
			}
			_, err := service.ListSeats(context.Background(), param)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
