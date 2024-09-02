package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"ticket-booking/internal/app"
	"ticket-booking/internal/app/service/customers"

	// "ticket-booking/internal/app/service/clerk"
	"ticket-booking/internal/domain/common"
	"ticket-booking/internal/domain/ecpay"
	response "ticket-booking/internal/router/api/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Post Order
// @Tags Customer
// @version 1.0
// @produce application/json
// @param page formData string true "event"
// @Success 200 string string "success get seats"
// @Failure 404 string errcode.error "no_found_item"
// @Failure 400 string errcode.error "invalid_parameter"
// @Router /api/v1/Seats/seat/ [get]
func CreateOrder(app *app.Application) gin.HandlerFunc {

	type Ticket struct {
		SeatNumber int    `json:"number"`
		Section    string `json:"string"`
		SeatType   string `json:"type"`
		Id         int32  `json:"ticket_id"`
		EventId    int32  `json:"event_id"`
	}

	type Body struct {
		ConsumerID int32   `json:"consumer_id" form:"consumer_id" binding:"required"`
		EventID    int32   `json:"event_id" form:"event_id" binding:"required"`
		SeatNumber []int32 `json:"seat_number" form:"seat_number" binding:"required"`
		// SeatAmount int32  `json:"seat_amount" form:"seat_amount" binding:"required"`
		Section string `json:"section" form:"section" binding:"required"`
	}

	type Response struct {
		Ticket        []Ticket `json:"tickets"`
		CheckMacValue string   `json:"check_mac_value"`
		OrderTradeNo  string   `json:"order_trade_no"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var body Body
		err := c.ShouldBind(&body)
		if err != nil {
			log.Panicf(err.Error())
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}
		orderTradeNo := uuid.New().String()

		seats, err := app.CustomerService.GetSeatsList(ctx,
			customers.SeatsParam{EventID: body.EventID, Section: body.Section, SeatNumberList: body.SeatNumber})
		// seat, err := app.CustomerService.GetSeat(ctx, customers.OrederParam{SeatID: body.SeatID})
		if err != nil {
			fmt.Println(err.Error())
			msg := "get seat fail"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}
		seatsIDlist := []int32{}
		for i := range seats {
			seatsIDlist = append(seatsIDlist, seats[i].ID)
		}

		ticket, err := app.CustomerService.CreateOrder(ctx, customers.OrederParam{
			ConsumerID:   body.ConsumerID,
			EventID:      body.EventID,
			SeatID:       seatsIDlist,
			OrderTradeNo: orderTradeNo,
			TicketStatus: 0,
		})
		if err != nil {
			if len(seatsIDlist) != 0 {
				app.CustomerService.UpdateSeatsAvailableBySeatID(ctx, customers.UpdateSeatsAvailableParm{
					SeatID:     seatsIDlist,
					SeatStatus: true,
				})
			}
			// log.Panicf(err.Error())
			msg := "create order fail"
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeResourceNotFound, errors.New(msg), common.WithMsg(msg)))
			return
		}
		data := ecpay.PaymentData{MerchantID: "3002607", MerchantTradeNo: "",
			MerchantTradeDate: "", PaymentType: "", TotalAmount: "", TradeDesc: "",
			ItemName: "", ReturnURL: "localhost/customer/paymented", ChoosePayment: ""}
		CheckMacValue := ecpay.GenerateCheckMacValue(data, "pwFHCqoQZGmho4w6", "EkRm7iFT261dpevs")
		// response result
		// ticketResp := ticketmanage.tieckResponse(ticket)
		resp := Response{
			CheckMacValue: CheckMacValue,
			OrderTradeNo:  orderTradeNo,
		}
		for i := range ticket {
			seat := seats[i]
			resp.Ticket = append(resp.Ticket, Ticket{
				SeatNumber: int(seat.SeatNumber),
				Section:    seat.Section,
				SeatType:   seat.Types,
				Id:         ticket[i].ID,
				EventId:    ticket[i].EventID,
			})
		}
		response.RespondWithJSON(c, http.StatusOK, resp)
	}
}

// @Summary Post Payment
// @Tags Customer
// @version 1.0
// @produce application/json
// @param page formData string true "event"
// @Success 200 string string "success get seats"
// @Failure 404 string errcode.error "no_found_item"
// @Failure 400 string errcode.error "invalid_parameter"
// @Router /api/v1/Seats/seat/ [get]
func Payment(app *app.Application) gin.HandlerFunc {

	type Body struct {
		OrderTradeNo  string `json:"order_trade_no"`
		CheckMacValue string `json:"check_mac_value"`
		RtnCode       int    `json:"rtn_code"`
	}
	type Response struct {
		OrderSuccess  bool   `json:"order_success"`
		CheckMacValue string `json:"check_mac_value"`
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var body Body
		err := c.ShouldBind(&body)
		OrderStaus := false
		if err != nil {
			log.Panicf(err.Error())
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
			return
		}

		// existm, err := app.CustomerService.CheckOrderExist(ctx, body.OrderTradeNo)
		// if err != nil {
		// 	msg := "fail to get OrderTradeNo"
		// 	response.RespondWithError(c,
		// 		common.NewError(common.ErrorCodeInternalProcess, errors.New(msg), common.WithMsg(msg)))
		// 	return
		// }
		// if !existm {
		// 	resp := Response{}
		// 	response.RespondWithJSON(c, http.StatusOK, resp)
		// 	return
		// }
		// app.CustomerService.CheckOrderExist(ctx)
		if body.RtnCode == 1 {
			_, err = app.CustomerService.UpdateTicketStatusIsPay(ctx, customers.UpdateTicketParm{
				OrderTradeNoParam: body.OrderTradeNo,
				TicketStatusParam: 1,
			})
			if err == nil {
				OrderStaus = true
			} else {
				msg := "fail to update ticket"
				response.RespondWithError(c,
					common.NewError(common.ErrorCodeInternalProcess, errors.New(msg), common.WithMsg(msg)))
				return
			}
		} else {
			err = app.CustomerService.UpdateTicketStatusTradeNoNotPay(ctx, customers.UpdateSeatsOrderTradeNoAvailableParm{
				OrderTradeNoParm: body.OrderTradeNo,
				SeatStatus:       true,
			})
			log.Panicf(err.Error())
			response.RespondWithError(c,
				common.NewError(common.ErrorCodeParameterInvalid, err, common.WithMsg("invalid parameter")))
		}

		resp := Response{
			OrderSuccess:  OrderStaus,
			CheckMacValue: body.OrderTradeNo,
		}
		response.RespondWithJSON(c, http.StatusOK, resp)
	}
}

func OrderDeadline(app *app.Application) {
	app.CustomerService.CheckOrderDeadlineAndCancel(context.Background())
}

func OrderStatusCheck(app *app.Application) {

	// defer cancel()
	app.CustomerService.CheckOrderExist(context.Background())
	// fmt.Println(err)

}
