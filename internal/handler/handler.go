package handler

import (
	"fmt"
	"net/http"

	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Handler is the main structure to inject dependencies.
type Handler struct {
	Log logger.Logger
	DB  db.DB
	EX  currency.Exchanger
}

// HealtResp holds the satus response from Health Check
type HealthResp struct {
	Status bool `json:"status"`
}

// Health method http GET
// @Summary Health check
// @Description Healthcheck endpoint, to ensure that the service is running.
// @Tags Health
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthResp
// @Router /health [get].
func (h Handler) Health(c *gin.Context) {
	if err := h.DB.Health(); err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		h.Log.Error(err)
		c.JSON(http.StatusInternalServerError, HealthResp{Status: false})
		return
	}

	c.JSON(http.StatusOK, &ResponseSuccess{&HealthResp{true}})
}
