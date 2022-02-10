package http

// ResponseSuccess describes an generic API response for success
type ResponseSuccess struct {
	Data interface{} `json:"data"`
}
