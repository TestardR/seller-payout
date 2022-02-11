package handler

import (

	// swagger docs.
	_ "github.com/TestardR/seller-payout/docs"
	"github.com/gin-gonic/gin"

	// swagger embed files.
	swaggerFiles "github.com/swaggo/files"

	// gin-swagger middleware.
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	healthRoute        = "/health"
	createItemsRoute   = "/items"
	readPayoutsRoute   = "/payouts/:seller_id"
	createSellersRoute = "/seller"
)

// @title SellerPayout Rest Server
// @version 1.0
// @description Server allowing interaction with Seller Payout Domain

// @contact.name Romain Testard
// @contact.email romain.rtestard@gmail.com

// @host localhost:3000

// NewServer instantiates an HTTP server.
func (h Handler) NewServer(env string) *gin.Engine {
	gin.SetMode(env)

	router := gin.New()
	router.Use(gin.Recovery())

	// swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health
	router.GET(healthRoute, h.Health)

	// Payouts
	router.GET(readPayoutsRoute, h.ReadPayouts)

	// Items
	router.POST(createItemsRoute, h.CreateItems)

	// Sellers
	router.POST(createSellersRoute, h.CreateSeller)

	return router
}
