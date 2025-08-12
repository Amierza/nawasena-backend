package response

import "time"

// Success response template for Swagger
type SwaggerResponseSuccess[T interface{}] struct {
	Status    bool        `json:"status" example:"true"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Data      T           `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
}

// Error response template for Swagger
type SwaggerResponseError struct {
	Status    bool        `json:"status" example:"false"`
	Message   string      `json:"message" example:"Invalid input"`
	Timestamp time.Time   `json:"timestamp"`
	Error     interface{} `json:"error"`
}
