package customers

import (
	"context"
	"database/sql"
	"errors"
	RedisCache "ticket-booking/internal/adapter/redisClient"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
	"ticket-booking/internal/domain/common"
	"ticket-booking/internal/domain/seatsmanage"
	"time"
)

type OrederParam struct {
	ConsumerID   int32
	EventID      int32
	SeatID       []int32
	TicketStatus int32
	OrderTradeNo string
}

func (c *CustomerService) CreateOrder(ctx context.Context, param OrederParam) ([]PostgresDB.Ticket, error) {

	ticket, err := c.orderRepo.CreateOrderAndSetSeatIsSoldTx(ctx, PostgresDB.CreateOrderAndSetSeatIsSoldParams{
		ConsumerID:   param.ConsumerID,
		EventID:      param.EventID,
		SeatID:       param.SeatID,
		OrderTradeNo: param.OrderTradeNo,
		TicketStatus: param.TicketStatus,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ticket, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to insert order")
		return ticket, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	// c.OrderMQ.
	err = c.orderMQ.SendOrderNotPay(ctx, param.OrderTradeNo)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to insert pay queue")
		return ticket, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	return ticket, err
}

type SeatsParam struct {
	EventID        int32
	Section        string
	SeatNumberList []int32
}

func (c *CustomerService) GetSeatsList(ctx context.Context, param SeatsParam) ([]seatsmanage.Seats, error) {
	result := []seatsmanage.Seats{}
	available, err := c.customerRedis.SetSeatKey(ctx, RedisCache.SeatKeyParams{
		EventID:    param.EventID,
		Section:    param.Section,
		SeatNumber: param.SeatNumberList,
	})
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get seat key")
		return result, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	// 加鎖後要釋出 上面err代表本身加鎖失敗
	defer func() {
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			c.customerRedis.DelSeatKey(ctx, RedisCache.SeatKeyParams{EventID: param.EventID,
				Section:    param.Section,
				SeatNumber: param.SeatNumberList})
		}
	}()

	if !available {
		if errors.Is(err, sql.ErrNoRows) {
			return result, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
	}

	seat, err := c.orderRepo.SelectSeatsListAndUpdateSeatStatusTx(ctx, PostgresDB.SelectSeatsListByEventAndSectionParams{
		EventID: param.EventID,
		Section: param.Section,
		Column3: param.SeatNumberList,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		// log.Panic(err)
		c.logger(ctx).Error().Err(err).Msg("failed to SelectSeatsListAndUpdateSeatStatusTx")
		return result, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	delcache := []PostgresDB.Seat{}
	for i := range seat {
		delitem := seat[i]
		delitem.SeatStatus = true
		delcache = append(delcache, delitem)
	}

	err = c.customerRedis.DelSeatValue(ctx, delcache)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to delete seat redis")
		return result, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	for i := range seat {
		result = append(result, seatsmanage.Seats{
			SeatStatus: seat[i].SeatStatus,
			Section:    seat[i].Section,
			Price:      seat[i].Price,
			ID:         seat[i].ID,
			SeatNumber: seat[i].SeatNumber,
		})
	}

	return result, err
}

type UpdateTicketParm struct {
	OrderTradeNoParam string
	TicketStatusParam int
}

func (c *CustomerService) UpdateTicketStatusIsPay(ctx context.Context, param UpdateTicketParm) ([]PostgresDB.UpdateTicketStatusRow, error) {
	tickets, err := c.orderRepo.UpdateTicketStatus(ctx, PostgresDB.UpdateTicketStatusParams{
		OrderTradeNo: param.OrderTradeNoParam,
		TicketStatus: int32(param.TicketStatusParam),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tickets, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return tickets, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	err = c.customerRedis.SetOrderTradeNo(ctx, param.OrderTradeNoParam)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to set OrderTradeNo")
	}
	return tickets, err
}

type UpdateSeatsOrderTradeNoAvailableParm struct {
	OrderTradeNoParm string
	SeatStatus       bool
}

func (c *CustomerService) UpdateTicketStatusTradeNoNotPay(ctx context.Context, param UpdateSeatsOrderTradeNoAvailableParm) error {
	seats, err := c.orderRepo.DeletTicketAndUpdateSeatIsAvailableStatusTx(ctx, param.OrderTradeNoParm, param.SeatStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return common.NewError(common.ErrorCodeInternalProcess, err)
	}
	SeatNumberList := []int32{}
	for i := range seats {
		SeatNumberList = append(SeatNumberList, seats[i].SeatNumber)
	}
	eventID := seats[0].EventID
	section := seats[0].Section
	err = c.customerRedis.DelSeatKey(ctx, RedisCache.SeatKeyParams{EventID: eventID,
		Section:    section,
		SeatNumber: SeatNumberList})
	if err != nil {
		return err
	}
	err = c.customerRedis.SetSeatAndIncrSeatNumber(ctx, []string{section}, seats)
	if err != nil {
		return err
	}
	err = c.customerRedis.SetOrderTradeNo(ctx, param.OrderTradeNoParm)
	if err != nil {
		return err
	}
	return err
}

type UpdateSeatsAvailableParm struct {
	SeatID     []int32
	SeatStatus bool
}

func (c *CustomerService) UpdateSeatsAvailableBySeatID(ctx context.Context, param UpdateSeatsAvailableParm) error {
	err := c.orderRepo.UpdateSeatStatusBySeatIDs(ctx, PostgresDB.UpdateSeatStatusBySeatIDsParams{
		SeatStatus: param.SeatStatus,
		Column2:    param.SeatID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return common.NewError(common.ErrorCodeInternalProcess, err)
	}
	return err
}
func (c *CustomerService) CheckOrderExist(ctx context.Context) {
	// Exist, err := c.orderMQ.ReceiverOrder(ctx, orderTradeNo)

	// var OrderMap chan map[string]int
	// var PaidOrder chan []string

	OrderMap := make(chan map[string]int)
	PaidOrder := make(chan []string)
	go c.orderMQ.ReceiverOrder(ctx, OrderMap, PaidOrder)

	IntervalTime := 3 * time.Second        // 觸發間隔時間
	ticker := time.NewTicker(IntervalTime) // 設定 2 秒觸發一次
	defer ticker.Stop()

	for range ticker.C {
		OrderTradeNoList, err := c.customerRedis.GetOrderTradeNoList(ctx)
		if err != nil {
			c.logger(ctx).Error().Err(err).Msg("failed to get orderTradeNo from redis")
		}
		OrderTradeNoMap := map[string]int{}
		for i := range OrderTradeNoList {
			OrderTradeNoMap[OrderTradeNoList[i]] = 1
		}
		if len(OrderTradeNoMap) > 0 {
			OrderMap <- OrderTradeNoMap
			select {
			case PaidOrderList := <-PaidOrder:
				if len(PaidOrderList) > 0 {
					c.customerRedis.DelOrderTradeNo(ctx, PaidOrderList)
				}
			default:
				break
			}
		}
	}
}
func (c *CustomerService) CheckOrderDeadlineAndCancel(ctx context.Context) error {
	ch := make(chan string, 200)
	err := c.orderMQ.DlxOrder(ctx, ch)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get orderTradeNo from MQ")
		return err
	}
	sem := make(chan struct{}, 100)
	go func() {
		sem <- struct{}{}
		for orderTradeNo := range ch {
			currentOrderTradeNo := orderTradeNo
			// 處理消息
			// fmt.Println("currentOrderTradeNo", currentOrderTradeNo)
			go func(orderTradeNo string) {
				defer func() { <-sem }() // release the semaphore slot
				// fmt.Println("currentOrderTradeNo", currentOrderTradeNo)
				exist, err := c.customerRedis.GetOrderTradeNoExist(ctx, orderTradeNo)
				if err != nil {
					c.logger(ctx).Error().Err(err).Msg("failed to  orderTradeNo DLX from redis")
				}
				if exist {
					return
				}
				// fmt.Println("dlx orderTradeNo", orderTradeNo)
				seats, err := c.orderRepo.DeletTicketAndUpdateSeatIsAvailableStatusTx(ctx, orderTradeNo, true)
				if err != nil {
					c.logger(ctx).Error().Err(err).Msg("failed to update orderTradeNo DLX from MQ")
				}
				if len(seats) == 0 {
					return
				}
				SeatNumberList := []int32{}
				for i := range seats {
					SeatNumberList = append(SeatNumberList, seats[i].SeatNumber)
				}
				eventID := seats[0].EventID
				section := seats[0].Section
				err = c.customerRedis.SetSeatAndIncrSeatNumber(ctx, []string{section}, seats)
				if err != nil {
					c.logger(ctx).Error().Err(err).Msg("failed to set seat value for orderTradeNo DLX")
				}
				err = c.customerRedis.DelSeatKey(ctx, RedisCache.SeatKeyParams{EventID: eventID,
					Section:    section,
					SeatNumber: SeatNumberList})
				// SeatNumber: SeatNumberList}) <<--- 這個可能有問題 沒有全部刪除
				if err != nil {
					c.logger(ctx).Error().Err(err).Msg("failed to delete seat key for orderTradeNo DLX")
				}
			}(currentOrderTradeNo)
		}
	}()
	return err
}
