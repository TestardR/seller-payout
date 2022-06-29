package cron

import (
	"errors"
	"testing"

	"github.com/TestardR/seller-payout/pkg/db"
	"github.com/TestardR/seller-payout/pkg/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_UpdateCurrencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mLog := mock.NewMockLogger(ctrl)
	mDB := mock.NewMockDB(ctrl)
	mEX := mock.NewMockExchanger(ctrl)

	h := handler{
		Log: mLog,
		DB:  mDB,
		EX:  mEX,
	}

	t.Run("should_be_ok", func(t *testing.T) {
		mLog.EXPECT().Info(gomock.Any())
		mDB.EXPECT().FindAll(gomock.Any())
		mDB.EXPECT().Update(gomock.Any())
		mLog.EXPECT().Info(gomock.Any())

		err := h.UpdateCurrencies()
		require.NoError(t, err)
	})

	t.Run("should_return_an_error_if_db_find_all_fails", func(t *testing.T) {
		mLog.EXPECT().Info(gomock.Any())
		mDB.EXPECT().FindAll(gomock.Any()).Return(errors.New("mock"))
		mLog.EXPECT().Error(gomock.Any())

		err := h.UpdateCurrencies()
		assert.ErrorIs(t, err, db.ErrDB)
	})

	t.Run("should_return_an_error_if_db_update_fails", func(t *testing.T) {
		mLog.EXPECT().Info(gomock.Any())
		mDB.EXPECT().FindAll(gomock.Any())
		mDB.EXPECT().Update(gomock.Any()).Return(errors.New("mock"))
		mLog.EXPECT().Error(gomock.Any())

		err := h.UpdateCurrencies()
		assert.ErrorIs(t, err, db.ErrDB)
	})
}
