package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type handlerCaseCreateSeller struct {
	h      Handler
	in     string
	status int
}

func TestHandler_CreateSeller(t *testing.T) {
	mc := gomock.NewController(t)
	t.Cleanup(func() { mc.Finish() })

	tests := map[string]handlerCaseCreateSeller{
		"fail-json":             sellerCreateCaseFailJSON(mc),
		"fail-empty-payload":    sellerCreateCaseFailEmptyPayload(mc),
		"fail-db-insert-seller": sellerCreateCaseFailDBInsertSeller(mc),
		"success":               sellerCreateCaseOK(mc),
	}

	for tn, tc := range tests {
		tn, tc := tn, tc
		t.Run(tn, func(t *testing.T) {
			router := NewServer(gin.TestMode, tc.h.Log, tc.h.DB)
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, createSellersRoute, bytes.NewBuffer([]byte(tc.in)))
			router.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.status {
				t.Errorf("Expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func sellerCreateCaseFailDBInsertSeller(mc *gomock.Controller) handlerCaseCreateSeller {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().Insert(gomock.Any()).Return(errors.New("mock"))
	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateSeller{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputSeller(),
		status: http.StatusInternalServerError,
	}
}

func sellerCreateCaseFailEmptyPayload(mc *gomock.Controller) handlerCaseCreateSeller {
	ml := mock.NewMockLogger(mc)

	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateSeller{
		h: Handler{
			Log: ml,
		},
		in:     "{}",
		status: http.StatusBadRequest,
	}
}

func sellerCreateCaseFailJSON(mc *gomock.Controller) handlerCaseCreateSeller {
	ml := mock.NewMockLogger(mc)

	ml.EXPECT().Error(gomock.Any())

	return handlerCaseCreateSeller{
		h: Handler{
			Log: ml,
		},
		in:     "{",
		status: http.StatusBadRequest,
	}
}

func sellerCreateCaseOK(mc *gomock.Controller) handlerCaseCreateSeller {
	ml := mock.NewMockLogger(mc)
	mdb := mock.NewMockDB(mc)

	mdb.EXPECT().Insert(gomock.Any())
	ml.EXPECT().Info(gomock.Any())

	return handlerCaseCreateSeller{
		h: Handler{
			Log: ml,
			DB:  mdb,
		},
		in:     validInputSeller(),
		status: http.StatusOK,
	}
}

func validInputSeller() string {
	return `{
		"currency": "EUR"
	}`
}
