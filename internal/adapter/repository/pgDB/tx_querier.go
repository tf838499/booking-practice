package PostgresDB

import (
	"context"
	"database/sql"
)

type CreateOrderAndSetSeatIsSoldParams struct {
	ConsumerID   int32
	EventID      int32
	SeatID       []int32
	TicketStatus int32
	OrderTradeNo string
}

func (store *PostgresRepository) CreateOrderAndSetSeatIsSoldTx(ctx context.Context, arg CreateOrderAndSetSeatIsSoldParams) ([]Ticket, error) {
	var result []Ticket
	var data Ticket
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		for i := range arg.SeatID {
			data, err = q.CreateOrder(ctx, CreateOrderParams{
				ConsumerID:   arg.ConsumerID,
				EventID:      arg.EventID,
				SeatID:       arg.SeatID[i],
				TicketStatus: arg.TicketStatus,
				OrderTradeNo: arg.OrderTradeNo,
			})
			result = append(result, data)
		}

		if err != nil {
			return err
		}
		return err
	})

	return result, err
}
func (store *PostgresRepository) SelectSeatsListAndUpdateSeatStatusTx(ctx context.Context, arg SelectSeatsListByEventAndSectionParams) ([]Seat, error) {

	var data []Seat
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		data, err = q.SelectSeatsListByEventAndSection(ctx, arg)
		if err != nil {
			return err
		}
		for i := range data {
			if !data[i].SeatStatus {
				return sql.ErrNoRows
			}
			row, err := q.UpdateSeatStatus(ctx, UpdateSeatStatusParams{ID: data[i].ID, SeatStatus: false})
			if err != nil {
				return err
			}
			data[i].SeatStatus = row.SeatStatus
		}
		return err
	})
	return data, err
}

func (store *PostgresRepository) DeletTicketAndUpdateSeatIsAvailableStatusTx(ctx context.Context, orderTradeNo string, seatStaus bool) ([]Seat, error) {

	var data []Seat
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		data, err = q.UpdateSeatStatusByOrderTradeNo(ctx, UpdateSeatStatusByOrderTradeNoParams{
			SeatStatus:   true,
			OrderTradeNo: orderTradeNo,
		})
		if err != nil {
			return err
		}
		err = q.DeleteTicketsByOrderTradeNo(ctx, orderTradeNo)
		return err
	})
	return data, err
}
