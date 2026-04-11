package middleware

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logging returns a request logging middleware using fiberzap.
func Logging(logger *zap.Logger) fiber.Handler {
	return fiberzap.New(fiberzap.Config{
		Logger: logger,
	})
}
