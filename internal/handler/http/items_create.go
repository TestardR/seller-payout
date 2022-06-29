package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

const (
	successMessage   = "success"
	numberOfDecimals = 3
)

var errMissingPayload = errors.New("there should be at least one item")

// CreateItemsRequest is the payload sent on the endpoint.
type CreateItemsRequest struct {
	Items []Item `validate:"dive"`
}

// Item is the payload expected out of the CreateItemsRequest.
type Item struct {
	Name     string    `json:"name" validate:"required"`
	Currency string    `required:"true" validate:"eq=GBP|eq=USD|eq=EUR"`
	Amount   int64     `json:"amount" validate:"required,min=0"`
	SellerID uuid.UUID `json:"seller_id" validate:"required"`
}

// CreateItems method http POST
// @Summary Endpoint to send sold items.
// @Description Create items.
// @Tags Items
// @Accept  json
// @Produce  json
// @Param create body handler.CreateItemsRequest true "Find the fields needed to create items using the 'handler' tab below."
// @Success 200 {object} ResponseSuccess
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /items [post].
func (h Handler) CreateItems(c *gin.Context) {
	outErr := func(status int, err error) {
		h.Log.Error(err)

		c.Error(err)
		c.JSON(status, newResponseError(err))
	}

	var input []Item
	if err := c.BindJSON(&input); err != nil {
		outErr(http.StatusBadRequest, fmt.Errorf("%w: %s", errBindJSON, err))

		return
	}

	if len(input) == 0 {
		outErr(http.StatusBadRequest, errMissingPayload)

		return
	}

	var req CreateItemsRequest
	req.Items = input

	if err := validator.New().Struct(req); err != nil {
		outErr(http.StatusBadRequest, fmt.Errorf("%w: %s", errValidatePayload, err))

		return
	}

	items, err := h.itemsFromInput(req.Items)
	if err != nil {
		outErr(http.StatusInternalServerError, err)

		return
	}

	if err := h.DB.Insert(&items); err != nil {
		outErr(http.StatusInternalServerError, fmt.Errorf("%w: %s", errDB, err))

		return
	}

	h.Log.Info(successMessage)
	c.JSON(http.StatusOK, &ResponseSuccess{items})
}

func (h Handler) itemsFromInput(input []Item) ([]domain.Item, error) {
	itemsDB := make([]domain.Item, 0, len(input))
	sellerMap := make(map[uuid.UUID]domain.Seller)

	// Note: if seller does not exist, we auto-create sellers with USD as currency for development sake
	// Not a good practice, in production, would get sellers through API or DB and discard unknown sellers.
	retrieveOrCreateSeller := func(item Item, sellerMap map[uuid.UUID]domain.Seller) (domain.Seller, error) {
		// cache seller to avoid unnecessary call.
		if s, ok := sellerMap[item.SellerID]; ok {
			return s, nil
		}

		var seller domain.Seller

		err := h.DB.FindByID(&seller, item.SellerID.String())
		if errors.Is(err, db.ErrRecordNotFound) {
			s := domain.Seller{ID: item.SellerID, CurrencyCode: currency.USDCode}
			sellerMap[item.SellerID] = s

			err := h.DB.Insert(&s)
			if err != nil {
				return domain.Seller{}, fmt.Errorf("%w: %s", errDB, err)
			}

			return sellerMap[item.SellerID], nil
		}

		if err != nil {
			return domain.Seller{}, fmt.Errorf("%w: %s", errDB, err)
		}

		sellerMap[item.SellerID] = seller

		return sellerMap[item.SellerID], nil
	}

	for _, item := range input {
		seller, err := retrieveOrCreateSeller(item, sellerMap)
		if err != nil {
			return nil, err
		}

		itemDB := domain.Item{
			ReferenceName: item.Name,
			Seller:        seller,
			SellerID:      item.SellerID,
			CurrencyCode:  item.Currency,
			PriceAmount:   decimal.NewFromInt(item.Amount).Round(numberOfDecimals),
		}

		itemsDB = append(itemsDB, itemDB)
	}

	return itemsDB, nil
}
