package seats

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	RedisCache "ticket-booking/internal/adapter/redisClient"
	PostgresDB "ticket-booking/internal/adapter/repository/pgDB"
	"ticket-booking/internal/domain/common"
	"ticket-booking/internal/domain/seatsmanage"
	"time"
)

type SectionParam struct {
	EventID int32
}

func (c *SeatService) SectionNumbers(ctx context.Context, param SectionParam) ([]seatsmanage.Section, error) {

	var EventID int32 = param.EventID
	SectionAndPrice, err := c.seatRedis.GetSectionsAndPricesByEventID(ctx, EventID)
	if err != nil {
		SectionAndPrice, err = c.seatRepo.SelectSectionsAndPricesByEventID(ctx, EventID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
			}
			c.logger(ctx).Error().Err(err).Msg("failed to get seat")
			return nil, common.NewError(common.ErrorCodeInternalProcess, err)
		}
	}
	sectionName := []string{}
	prices := []float64{}
	for i := range SectionAndPrice {
		sectionName = append(sectionName, SectionAndPrice[i].Section)
		prices = append(prices, SectionAndPrice[i].Price)
	}
	count, err := c.seatRedis.GetSectionSeatCount(ctx, RedisCache.GetSectionSeatCount{EventId: int(EventID), Section: sectionName})
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get section amount")
		return nil, err
	}
	sections := seatsmanage.ArrageSectionNumbersPrices(sectionName, prices, count)
	return sections, err
}

type SeatListParam struct {
	EventID int32
	Section string
}

func (c *SeatService) ListSeats(ctx context.Context, param SeatListParam) ([]PostgresDB.Seat, error) {

	var EventID int32 = param.EventID
	var Section string = param.Section
	SectionPrice, err := c.seatRedis.GetSectionsPricesBySectionsName(ctx, EventID, Section)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get section price form redis")
		return nil, common.NewError(common.ErrorCodeResourceNotFound, sql.ErrNoRows)
	}

	sections := []string{Section}
	prices := []float64{SectionPrice}

	seats, err := c.seatRedis.GetSeatValue(ctx, RedisCache.GetSeatParams{EventId: int(EventID), Section: sections, Price: prices})

	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get good form redis")
	} else if len(seats) != 0 {
		seats = seatsmanage.SortSeatsBySectionAndSeatNumber(seats)
		return seats, err
	} else if len(seats) == 0 {
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return nil, common.NewError(common.ErrorCodeResourceNotFound, sql.ErrNoRows)
	}

	return seats, err
}
func (c *SeatService) PreheatListSeats(ctx context.Context, param SeatListParam) ([]PostgresDB.Seat, error) {

	var EventID int32 = param.EventID

	SectionAndPrice, err := c.seatRepo.SelectSectionsAndPricesByEventID(ctx, EventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	sections := []string{}
	prices := []float64{}
	for i := range SectionAndPrice {
		sections = append(sections, SectionAndPrice[i].Section)
		prices = append(prices, SectionAndPrice[i].Price)
	}
	seats, err := c.seatRedis.GetSeatValue(ctx, RedisCache.GetSeatParams{EventId: int(EventID), Section: sections, Price: prices})
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to get good form redis")
	} else if len(seats) != 0 {
		seats = seatsmanage.SortSeatsBySectionAndSeatNumber(seats)
		return seats, err
	}
	// data := c.seatRedis.GetSectionSeatAmount(ctx, RedisCache.GetSectionSeatAmount{EventId: int(EventID), Section: sections})
	seats, err = c.seatRepo.SelectSeatsList(ctx, EventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		c.logger(ctx).Error().Err(err).Msg("failed to get seat")
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	err = c.seatRedis.SetSeatValueAndSectionAndIncrSeatNumber(ctx, sections, prices, seats)
	if err != nil {
		c.logger(ctx).Error().Err(err).Msg("failed to insert seat")
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	return seats, err
}

func (c *SeatService) CreateInitDBdata(ctx context.Context) error {
	eventdata := time.Now().Local()
	Eventdata, err := c.dataRepo.CreateEvent(ctx,
		PostgresDB.CreateEventParams{EventName: "Concert A", EventDate: eventdata, TotalSeats: 10000})
	if err != nil {
		log.Fatalf("Failed to create Event  %v", err)
		return err
	}
	for i := 1; i <= 1000; i++ {
		consumer, err := c.dataRepo.CreateConsumer(ctx, PostgresDB.CreateConsumerParams{
			Name:  fmt.Sprintf("Consumer %d", i),
			Email: fmt.Sprintf("consumer%d@example.com", i),
		})
		if err != nil {
			log.Fatalf("Failed to create consumer %d: %v", i, err)
			return err
		}
		fmt.Printf("Created consumer: %+v\n", consumer)
	}

	// Create 1000 seats
	for section := 1; section <= 10; section++ {
		for seatNumber := 1; seatNumber <= 1000; seatNumber++ {
			seat, err := c.dataRepo.CreateSeat(ctx, PostgresDB.CreateSeatParams{
				Section:    fmt.Sprintf("Section %d", section),
				SeatNumber: int32(seatNumber),
				SeatStatus: true,
				Price:      50.0 * float64(section%10), // random price
				EventID:    Eventdata.ID,               // assigning event ID 1
			})
			if err != nil {
				log.Fatalf("Failed to create seat in section %d, seat number %d: %v", section, seatNumber, err)
			}
			fmt.Printf("Created seat: %+v\n", seat)
		}
	}
	// EventID: body.Event,
	// 			Section: []string{
	// 				"Section 1", "Section 2", "Section 3", "Section 4", "Section 5",
	// 				"Section 6", "Section 7", "Section 8", "Section 9", "Section 10",
	// 			}}
	c.seatRedis.FlushAll(ctx, RedisCache.DelSeatParams{
		EventId: 1,
		Section: []string{
			"Section 1", "Section 2", "Section 3", "Section 4", "Section 5",
			"Section 6", "Section 7", "Section 8", "Section 9", "Section 10",
		},
	})
	return err
}
