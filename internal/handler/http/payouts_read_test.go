package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/TestardR/seller-payout/pkg/currency"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type handlerCaseReadPayouts struct {
	h        Handler
	sellerID string
	status   int
}

func TestHandler_ReadPayouts(t *testing.T) {
	mc := gomock.NewController(t)
	t.Cleanup(func() { mc.Finish() })

	tests := map[string]handlerCaseReadPayouts{
		"fail-db-find-by-seller":   payoutsReadCaseFailDBFindSellerByID(mc),
		"fail-db-seller-not-found": payoutsReadCaseFailDBSellerNotFound(mc),
		"fail-db-find-payouts":     payoutsReadCaseFailDBFindPayouts(mc),
		"success":                  payoutsReadCaseOK(mc),
	}

	for tn, tc := range tests {
		tn, tc := tn, tc
		t.Run(tn, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)

			router.GET("/:seller_id", tc.h.ReadPayouts)
			uri := fmt.Sprintf("/%s", tc.sellerID)

			var err error
			ctx.Request, err = http.NewRequest("GET", uri, nil)
			require.NoError(t, err)

			router.ServeHTTP(w, ctx.Request)

			if w.Result().StatusCode != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func payoutsReadCaseOK(mc *gomock.Controller) handlerCaseReadPayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().FindByID(gomock.Any(), "123")
	mdb.EXPECT().FindPayoutsBySellerID("123").Return([]domain.Payout{}, nil)
	ml.EXPECT().Info(gomock.Any())

	return handlerCaseReadPayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		sellerID: "123",
		status:   http.StatusOK,
	}
}

func payoutsReadCaseFailDBFindPayouts(mc *gomock.Controller) handlerCaseReadPayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().FindByID(gomock.Any(), "123")
	mdb.EXPECT().FindPayoutsBySellerID("123").Return([]domain.Payout{}, errors.New("mock"))
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseReadPayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		sellerID: "123",
		status:   http.StatusInternalServerError,
	}
}

func payoutsReadCaseFailDBSellerNotFound(mc *gomock.Controller) handlerCaseReadPayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().FindByID(gomock.Any(), "123").Return(db.ErrRecordNotFound)
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseReadPayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		sellerID: "123",
		status:   http.StatusBadRequest,
	}
}

func payoutsReadCaseFailDBFindSellerByID(mc *gomock.Controller) handlerCaseReadPayouts {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().FindByID(gomock.Any(), "123").Return(errors.New("mock"))
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseReadPayouts{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		sellerID: "123",
		status:   http.StatusInternalServerError,
	}
}

func Test_NewPayoutsFromInput(t *testing.T) {

	expected := []domain.Payout{
		{ID: uuid.FromStringOrNil("test"),
			PriceTotal: decimal.NewFromInt(1),
			CreatedAt:  time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			Currency:   domain.Currency{Code: currency.EURCode},
			Items: []domain.Item{
				{
					ID:            uuid.FromStringOrNil("test"),
					ReferenceName: "test",
				},
			},
		},
	}

	got := newPayoutsFromInput(expected)

	assert.Equal(t, got[0].ID, expected[0].ID)
	assert.Equal(t, got[0].Items[0].ID, expected[0].Items[0].ID)
}
