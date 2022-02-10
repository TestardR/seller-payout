package http

import (
	"net/http"
	"time"

	"github.com/TestardR/seller-payout/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AccessLog middleware for logging every request incomming into the server
func AccessLog(appLogger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()
		stop := time.Now()

		data := make(map[string]interface{})

		data["duration"] = stop.Sub(start)
		data["network.client.ip"] = c.ClientIP()
		data["network.bytes_written"] = c.Writer.Size()
		data["http.method"] = c.Request.Method
		data["http.url_details.path"] = path
		data["http.status_code"] = c.Writer.Status()

		if len(query) > 0 {
			data["http.url_details.queryString"] = query
		}

		// Write log with level according to status
		if c.Writer.Status() >= http.StatusInternalServerError {
			appLogger.Error(data)
		} else {
			appLogger.Info(data)
		}
	}
}
