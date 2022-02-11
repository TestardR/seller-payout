package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseError(t *testing.T) {
	err1 := errors.New("foo")
	err2 := errors.New("bar")

	respError := newResponseError(err1, err2)

	assert.Len(t, respError.Errs, 2)
	assert.Equal(t, "foo", respError.Errs[0].Message)
	assert.Equal(t, "bar", respError.Errs[1].Message)
}
