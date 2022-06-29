package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Health(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mLog := mock.NewMockLogger(ctrl)
	mDB := mock.NewMockDB(ctrl)
	h := handler{
		Log: mLog,
		DB:  mDB,
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET(healthRoute, h.Health)

	ts := httptest.NewServer(router)
	defer ts.Close()

	t.Run("should_be_ok", func(t *testing.T) {
		mDB.EXPECT().Health().Return(nil)
		resp, err := http.Get(fmt.Sprintf("%s%s", ts.URL, healthRoute))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should_return_500", func(t *testing.T) {
		mDB.EXPECT().Health().Return(errors.New("mock"))
		mLog.EXPECT().Error(gomock.Any())
		resp, err := http.Get(fmt.Sprintf("%s%s", ts.URL, healthRoute))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

}
