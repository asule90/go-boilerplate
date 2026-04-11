package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

// Common returns the standard middleware stack applied to all routes.
func Common(logger *zap.Logger) []fiber.Handler {
	return []fiber.Handler{
		requestid.New(),
		recover.New(),
		CORS(),
		Logging(logger),
	}
}
