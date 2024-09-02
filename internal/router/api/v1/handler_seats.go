package v1

import (
	"errors"
	"net/http"

	"ticket-booking/internal/app"
	"ticket-booking/internal/app/service/seats"

	// "ticket-booking/internal/app/service/clerk"
	// "ticket-booking/internal/domain/backstage"
	"ticket-booking/internal/domain/common"
	response "ticket-booking/internal/router/api/response"

	"github.com/gin-gonic/gin"
)

// @Summary Get Sections
// @Tags Seats
// @version 1.0
// @produce application/json
// @param page formData string true "event"
// @Success 200 string string "success get seats"
// @Failure 404 string errcode.error "no_found_item"
// @Failure 400 string errcode.error "invalid_parameter"
// @Router /api/v1/Seats/seat/ [get]
func Sections(app *app.Application) gin.HandlerFunc {

	type Section struct {
		Section   string  `json:"section"`
		SeatPrice float64 `json:"price"`
		Count     int32   `json:"count"`
	}

	type Body struct {
		Event int32 `form:"event" binding:"required"`
	}

	type Response struct {
		Sections []Section `json:"section"`
	}

	return func(c *gin.Context) {

		ctx := c.Request.Context()
		var body Body
		err := c.ShouldBindQuery(&body)
		if err != nil {
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}
		SectionList, err := app.SeatService.SectionNumbers(ctx,
			seats.SectionParam{
				EventID: body.Event,
			})
		if err != nil {
			msg := "no found section"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}
		resp := Response{Sections: []Section{}}
		for i := range SectionList {
			resp.Sections = append(resp.Sections, Section{
				Section:   SectionList[i].Section,
				SeatPrice: SectionList[i].Price,
				Count:     SectionList[i].Count,
			})
		}
		response.RespondWithJSON(c, http.StatusOK, resp)
	}
}

// @Summary Get Seats
// @Tags Seats
// @version 1.0
// @produce application/json
// @param page formData string true "event"
// @Success 200 string string "success get seats"
// @Failure 404 string errcode.error "no_found_item"
// @Failure 400 string errcode.error "invalid_parameter"
// @Router /api/v1/Seats/seat/ [get]
func ListSeats(app *app.Application) gin.HandlerFunc {

	type Seat struct {
		SeatNumber int32   `json:"number"`
		Section    string  `json:"section"`
		SeatPrice  float64 `json:"price"`
		SeatStatus bool    `json:"seat_status"`
	}

	type Body struct {
		Event   int32  `form:"event" binding:"required"`
		Section string `form:"section" binding:"required"`
	}

	type Response struct {
		Seats []Seat `json:"seats"`
	}

	return func(c *gin.Context) {

		ctx := c.Request.Context()
		var body Body
		err := c.ShouldBindQuery(&body)
		if err != nil {
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}

		// ListSeats(ctx, seats.SeatListParam{EventID: body.Event})
		seats, err := app.SeatService.ListSeats(ctx,
			seats.SeatListParam{
				EventID: body.Event,
				Section: body.Section,
			})
		if err != nil {
			msg := "no found item"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}

		resp := Response{Seats: []Seat{}}
		for i := range seats {
			seat := seats[i]
			resp.Seats = append(resp.Seats, Seat{
				Section:    seat.Section,
				SeatStatus: seat.SeatStatus,
				SeatNumber: seat.SeatNumber,
				SeatPrice:  seat.Price,
			})
		}
		response.RespondWithJSON(c, http.StatusOK, resp)

	}
}

// @Summary Get Seats
// @Tags Seats
// @version 1.0
// @produce application/json
// @param page formData string true "event"
// @Success 200 string string "success get seats"
// @Failure 404 string errcode.error "no_found_item"
// @Failure 400 string errcode.error "invalid_parameter"
// @Router /api/v1/Seats/seat/ [get]
func PreHeatListSeats(app *app.Application) gin.HandlerFunc {

	type Seat struct {
		SeatNumber int32   `json:"number"`
		Section    string  `json:"section"`
		SeatPrice  float64 `json:"price"`
		SeatStatus bool    `json:"seat_status"`
	}

	type Body struct {
		Event int32 `form:"event" binding:"required"`
	}

	type Response struct {
		Seats []Seat `json:"seats"`
	}

	return func(c *gin.Context) {

		ctx := c.Request.Context()
		var body Body
		err := c.ShouldBindQuery(&body)
		if err != nil {
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}
		// ListSeats(ctx, seats.SeatListParam{EventID: body.Event})
		seats, err := app.SeatService.PreheatListSeats(ctx,
			seats.SeatListParam{
				EventID: body.Event,
			})
		if err != nil {
			msg := "no found item"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}

		resp := Response{Seats: []Seat{}}
		for i := range seats {
			seat := seats[i]
			resp.Seats = append(resp.Seats, Seat{
				Section:    seat.Section,
				SeatStatus: seat.SeatStatus,
				SeatNumber: seat.SeatNumber,
				SeatPrice:  seat.Price,
			})
		}
		response.RespondWithJSON(c, http.StatusOK, resp)

	}
}
func CreatDBdata(app *app.Application) gin.HandlerFunc {

	type Body struct {
		Event int32 `form:"event" binding:"required"`
	}

	type Response struct {
		Data string `json:"resp"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		err := app.SeatService.CreateInitDBdata(ctx)
		if err != nil {
			msg := "fail create"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}
		resp := Response{Data: "success"}
		response.RespondWithJSON(c, http.StatusOK, resp)
	}
}
