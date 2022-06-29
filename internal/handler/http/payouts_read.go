package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type payout struct {
	ID        uuid.UUID       `json:"id"`
	Price     decimal.Decimal `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
	Currency  string          `json:"currency"`
	Items     []item          `json:"items"`
}

type item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// ReadPayouts method http GET
// @Summary Endpoint to retrieve payouts for a specific seller.
// @Description Create Seller.
// @Tags Seller
// @Accept  json
// @Produce  json
// @Param seller_id query string true "Seller ID query parameter"
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /payouts/:seller_id [get].
func (h handler) ReadPayouts(c *gin.Context) {
	outErr := func(status int, err error) {
		h.Log.Error(err)

		c.Error(err)
		c.JSON(status, newResponseError(err))
	}

	sellerID := c.Param("seller_id")

	var seller domain.Seller

	err := h.DB.FindByID(&seller, sellerID)
	if errors.Is(err, db.ErrRecordNotFound) {
		outErr(http.StatusBadRequest, fmt.Errorf("%w: %s", db.ErrDB, err))

		return
	}

	if err != nil {
		outErr(http.StatusInternalServerError, fmt.Errorf("%w: %s", db.ErrDB, err))

		return
	}

	payouts, err := h.DB.FindPayoutsBySellerID(sellerID)
	if err != nil {
		outErr(http.StatusInternalServerError, fmt.Errorf("%w: %s", db.ErrDB, err))

		return
	}

	p := newPayoutsFromInput(payouts)

	h.Log.Info(successMessage)
	c.JSON(http.StatusOK, &ResponseSuccess{p})
}

func newItemsFromInput(items []domain.Item) []item {
	output := make([]item, 0, len(items))

	for _, dbItem := range items {
		it := item{
			ID:   dbItem.ID,
			Name: dbItem.ReferenceName,
		}
		output = append(output, it)
	}

	return output
}

func newPayoutsFromInput(dbPayouts []domain.Payout) []payout {
	output := make([]payout, 0, len(dbPayouts))

	for _, DBpayout := range dbPayouts {
		p := payout{
			ID:        DBpayout.ID,
			Price:     DBpayout.PriceTotal,
			CreatedAt: DBpayout.CreatedAt,
			Currency:  DBpayout.Currency.Code,
			Items:     newItemsFromInput(DBpayout.Items),
		}

		output = append(output, p)
	}

	return output
}
