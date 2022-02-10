package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Health(t *testing.T) {
	router := NewServer(gin.TestMode, nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, healthRoute, nil)
	router.ServeHTTP(w, req)

	b, _ := json.Marshal(&ResponseSuccess{&HealthResp{true}})

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(b), w.Body.String())
}
