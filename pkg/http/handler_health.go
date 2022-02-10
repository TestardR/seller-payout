package http

import "github.com/gin-gonic/gin"

type HealthResp struct {
	Status bool `json:"status"`
}

func (h handler) Health(c *gin.Context) {
	c.JSON(200, &ResponseSuccess{&HealthResp{true}})
}
