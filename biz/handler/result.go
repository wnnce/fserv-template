// base.go
// Author:      cola
// Description: TODO: Describe this file
// Created:     2025/7/26 16:16

package handler

import (
	"github.com/gofiber/fiber/v3"
	"time"
)

// Result is a generic response wrapper used for standardizing API responses.
// It includes a status code, message, optional timestamp, and a data payload.
type Result[T any] struct {
	Code      int    `json:"code"`                // Business status code
	Message   string `json:"message,omitempty"`   // Message or error description
	Timestamp int64  `json:"timestamp,omitempty"` // Server response timestamp (milliseconds)
	Data      T      `json:"data,omitempty"`      // Response payload
}

// NewResult creates a new generic Result with the given code, message, and data.
// It automatically adds a Unix millisecond timestamp.
func NewResult[T any](code int, message string, data T) *Result[T] {
	return &Result[T]{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
}

// Ok returns a success Result with status code 200 and no data.
func Ok() *Result[any] {
	return NewResult[any](fiber.StatusOK, "ok", nil)
}

// OkWithData returns a success Result with status code 200 and a data payload.
func OkWithData[T any](data T) *Result[T] {
	return NewResult[T](fiber.StatusOK, "ok", data)
}

// Fail returns a failure Result with the given error code and message, and no data.
func Fail(code int, message string) *Result[any] {
	return NewResult[any](code, message, nil)
}

// FailWithRequest returns a failure Result for client-side errors (HTTP 400).
func FailWithRequest(message string) *Result[any] {
	return NewResult[any](fiber.StatusBadRequest, message, nil)
}

// FailWithServer returns a failure Result for internal server errors (HTTP 500).
func FailWithServer(message string) *Result[any] {
	return NewResult[any](fiber.StatusInternalServerError, message, nil)
}

// FailWithAuth returns a failure Result for authentication errors (HTTP 401).
func FailWithAuth(message string) *Result[any] {
	return NewResult[any](fiber.StatusUnauthorized, message, nil)
}
