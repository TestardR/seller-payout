package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TestardR/seller-payout/internal/domain"
	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type handlerCaseCreateItems struct {
	h      handler
	in     string
	status int
}

func TestHandler_CreateItems(t *testing.T) {
	mc := gomock.NewController(t)
	t.Cleanup(func() { mc.Finish() })

	tests := map[string]handlerCaseCreateItems{
		"fail-json":                  itemsCreateCaseFailJSON(mc),
		"fail-empty-payload":         itemsCreateCaseFailEmptyPayload(mc),
		"fail-validation":            itemsCreateCaseFailValidation(mc),
		"fail-db-find-seller-by-id":  itemsCreateCaseFailDBFindSellerByID(mc),
		"fail-db-insert-items":       itemsCreateCaseFailDBInsertItems(mc),
		"auto-create-seller-success": itemsCreateCaseAutoCreateSeller(mc),
		"success":                    itemsCreateCaseOK(mc),
	}

	for tn, tc := range tests {
		tn, tc := tn, tc
		t.Run(tn, func(t *testing.T) {
			router := NewServer(gin.TestMode, tc.h.Log, tc.h.DB)
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, createItemsRoute, bytes.NewBuffer([]byte(tc.in)))
			router.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func itemsCreateCaseFailEmptyPayload(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)
	mDB := mock.NewMockDB(mc)

	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
			DB:  mDB,
		},
		in:     "[]",
		status: http.StatusBadRequest,
	}
}

func itemsCreateCaseFailJSON(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)

	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
		},
		in: `[ {
        "name": "bag,
        "amount": 1,
        "seller_id": "78dd7916-f276-494b-84a8-83e5bbee8cc6"
    }]`,
		status: http.StatusBadRequest,
	}
}

func itemsCreateCaseFailValidation(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)

	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
		},
		in: `[ {
        "name": "bag",
        "amount": -1,
        "seller_id": "78dd7916-f276-494b-84a8-83e5bbee8cc6"
    }]`,
		status: http.StatusBadRequest,
	}
}

func itemsCreateCaseFailDBFindSellerByID(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().FindByID(&domain.Seller{}, "78dd7916-f276-494b-84a8-83e5bbee8c11").Return(errors.New("mock"))
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputItems(),
		status: http.StatusInternalServerError,
	}
}

func itemsCreateCaseFailDBInsertItems(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mSellerID := "78dd7916-f276-494b-84a8-83e5bbee8c11"

	mdb.EXPECT().FindByID(&domain.Seller{}, mSellerID)
	mdb.EXPECT().Insert(gomock.Any()).Return(errors.New("mock"))
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputItems(),
		status: http.StatusInternalServerError,
	}
}

func itemsCreateCaseAutoCreateSeller(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mSellerID := "78dd7916-f276-494b-84a8-83e5bbee8c11"

	mdb.EXPECT().FindByID(&domain.Seller{}, mSellerID).Return(db.ErrRecordNotFound)
	mdb.EXPECT().Insert(gomock.Any())
	mdb.EXPECT().Insert(gomock.Any())
	ml.EXPECT().Info(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputItems(),
		status: http.StatusOK,
	}
}

func itemsCreateCaseOK(mc *gomock.Controller) handlerCaseCreateItems {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mSellerID := "78dd7916-f276-494b-84a8-83e5bbee8c11"

	mdb.EXPECT().FindByID(&domain.Seller{}, mSellerID)
	mdb.EXPECT().Insert(gomock.Any())
	ml.EXPECT().Info(gomock.Any())

	return handlerCaseCreateItems{
		h: handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputItems(),
		status: http.StatusOK,
	}
}

func validInputItems() string {
	return `[
		{
			"name": "bag",
			"amount": 1,
			"currency": "GBP",
			"seller_id": "78dd7916-f276-494b-84a8-83e5bbee8c11"
		}
	]`
}
