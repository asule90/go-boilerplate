package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/internal/db"
	"github.com/sule/go-boilerplate/internal/middleware"
	"github.com/sule/go-boilerplate/internal/user"
	"go.uber.org/zap"
)

// New creates and configures the Fiber application with all routes and middleware.
func New(cfg *config.Provider, pool *pgxpool.Pool, logger *zap.Logger, fbAuth *middleware.FirebaseAuth) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: defaultErrorHandler(logger),
	})

	for _, m := range middleware.Common(logger) {
		app.Use(m)
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": cfg.App.Name})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	qb := db.NewQueryBuilder(pool, logger)

	userRepo := user.NewRepository(qb)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc, cfg, logger, fbAuth)

	api := app.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	return app
}

func defaultErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		logger.Error("unhandled fiber error", zap.Error(err))
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
}
