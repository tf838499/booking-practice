package router

import (
	"github.com/gin-gonic/gin"

	"ticket-booking/internal/app"
	v1 "ticket-booking/internal/router/api/v1"
)

func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerAPIHandlers(router, app)
}

func registerAPIHandlers(router *gin.Engine, app *app.Application) {
	// Build middlewares
	// BearerToken := NewAuthMiddlewareBearer(app)
	// OauthToken := NewOAuthMiddleware(app)
	// We mount all handlers under /api path
	r := router.Group("/api")
	v := r.Group("/v1")
	// v.GET("/callback", OauthToken.Callback)

	// Add barter namespace
	seatGroup := v.Group("/seat")
	{
		seatGroup.GET("/sections", v1.Sections(app))
		seatGroup.GET("/", v1.ListSeats(app))
		seatGroup.GET("/preheat", v1.PreHeatListSeats(app))
		seatGroup.GET("/initdb", v1.CreatDBdata(app))
	}
	customerGroup := v.Group("/customer")
	{
		customerGroup.POST("/order", v1.CreateOrder(app))
		customerGroup.POST("/ecpay/return", v1.Payment(app))
	}
	go v1.OrderDeadline(app)
	go v1.OrderStatusCheck(app)
}
