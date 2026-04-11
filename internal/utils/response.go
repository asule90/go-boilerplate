package utils

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sule/go-boilerplate/pkg/errr"
	"go.uber.org/zap"
)

// Envelope is the standard JSON response wrapper.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ResponseOK(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Envelope{Success: true, Data: data})
}

func ResponseCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(Envelope{Success: true, Data: data})
}

func ResponseOKWithMeta(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.Status(http.StatusOK).JSON(Envelope{Success: true, Data: data, Meta: meta})
}

func ResponseError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Envelope{Success: false, Message: message})
}

func ResponseValidationError(c *fiber.Ctx, errs interface{}) error {
	return c.Status(http.StatusBadRequest).JSON(Envelope{
		Success: false,
		Message: "validation failed",
		Errors:  errs,
	})
}

func ResponseNotFound(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusNotFound).JSON(Envelope{Success: false, Message: message})
}

func ResponseInternalError(c *fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(Envelope{
		Success: false,
		Message: "internal server error",
	})
}

// ParseErrHTTP maps an error to the appropriate HTTP response.
func ParseErrHTTP(c *fiber.Ctx, err error, payload interface{}, logger *zap.Logger) error {
	var sce *errr.StatusCodeError
	if errors.As(err, &sce) {
		return ResponseError(c, sce.Code, sce.Message)
	}

	if errors.Is(err, errr.ErrNoRows) {
		return ResponseNotFound(c, "resource not found")
	}

	if logger != nil {
		logger.Error("unhandled error", zap.Error(err))
	}
	return ResponseInternalError(c)
}
