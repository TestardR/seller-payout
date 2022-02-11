package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("should return an errParseEnv error", func(t *testing.T) {
		_, err := config()

		assert.Equal(t, true, errors.Is(err, errParseEnv))
	})

	t.Run("should return an errParseEnv error because validator failed", func(t *testing.T) {
		t.Setenv("PORT", "v")
		t.Setenv("ENV", "v")
		t.Setenv("PAYOUT_INTERVAL", "v")
		t.Setenv("CURRENCY_INTERVAL", "v")
		t.Setenv("PG_HOST", "v")

		_, err := config()
		assert.Equal(t, true, errors.Is(err, errParseEnv))
	})

	t.Run("should be ok", func(t *testing.T) {
		t.Setenv("PORT", "v")
		t.Setenv("ENV", "debug")
		t.Setenv("PAYOUT_INTERVAL", "4")
		t.Setenv("CURRENCY_INTERVAL", "12")
		t.Setenv("PG_HOST", "postgres")
		t.Setenv("PG_USER", "v")
		t.Setenv("PG_NAME", "postgres")
		t.Setenv("PG_PASSWORD", "v")

		_, err := config()
		require.NoError(t, err)
	})
}
