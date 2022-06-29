package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var (
	errBindJSON        = errors.New("failed to decode JSON payload")
	errValidatePayload = errors.New("failed to validate payload")
	errDB              = errors.New("failed to perform database operation")
)

// Seller is the item owner with a desired currency for payouts.
type Seller struct {
	Currency string `required:"true" validate:"eq=GBP|eq=USD|eq=EUR"`
}

// CreateSeller method http POST
// @Summary Endpoint to create seller.
// @Description Create Seller.
// @Tags Seller
// @Accept  json
// @Produce  json
// @Param create body handler.Seller true "Find the fields needed to create a seller using the 'handler' tab below."
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /seller [post].
func (h Handler) CreateSeller(c *gin.Context) {
	var input Seller

	outErr := func(status int, err error) {
		h.Log.Error(err)

		c.Error(err)
		c.JSON(status, newResponseError(err))
	}

	if err := c.BindJSON(&input); err != nil {
		outErr(http.StatusBadRequest, fmt.Errorf("%w: %s", errBindJSON, err))

		return
	}

	if err := validator.New().Struct(input); err != nil {
		outErr(http.StatusBadRequest, fmt.Errorf("%w: %s", errValidatePayload, err))

		return
	}

	seller := domain.Seller{
		CurrencyCode: input.Currency,
	}

	if err := h.DB.Insert(&seller); err != nil {
		outErr(http.StatusInternalServerError, fmt.Errorf("%w: %s", errDB, err))

		return
	}

	h.Log.Info(successMessage)
	c.JSON(http.StatusOK, &ResponseSuccess{seller})
}
