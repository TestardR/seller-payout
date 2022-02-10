package http

import (
	"github.com/TestardR/seller-payout/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	healthRoute = "/health"
)

type handler struct {
	log logger.Logger
}

func NewServer(env string, log logger.Logger) *gin.Engine {
	h := handler{log}

	gin.SetMode(env)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(AccessLog(log))

	// useful for monitoring our service and CI/CD tools
	router.GET(healthRoute, h.Health)

	return router
}
