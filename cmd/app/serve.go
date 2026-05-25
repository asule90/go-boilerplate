package app

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/internal/db"
	"github.com/sule/go-boilerplate/internal/middleware"
	"github.com/sule/go-boilerplate/internal/server"
	"github.com/sule/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE:  runServe,
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	log, err := logger.InitLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer log.Sync()

	log.Info("bootstrap configuration",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
		zap.String("timezone", cfg.App.Timezone),
		zap.String("auth_mode", cfg.Auth.Mode),
		zap.String("database", dbConnectionTarget(cfg.Database.URL)),
	)

	pool, err := db.NewPgxPool(context.Background(), cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	log.Info("database connection established", zap.String("database", dbConnectionTarget(cfg.Database.URL)))

	var fbAuth *middleware.FirebaseAuth
	if cfg.Auth.Mode == "firebase" {
		fbAuth, err = middleware.InitFirebase(cfg)
		if err != nil {
			return fmt.Errorf("failed to init firebase: %w", err)
		}
	}

	app := server.New(cfg, pool, log, fbAuth)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.App.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	<-quit
	return app.Shutdown()
}

func dbConnectionTarget(databaseURL string) string {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "unknown"
	}

	host := u.Hostname()
	port := u.Port()
	dbName := strings.TrimPrefix(u.Path, "/")

	if host == "" {
		host = "unknown-host"
	}
	if dbName == "" {
		dbName = "unknown-db"
	}

	if port != "" {
		return fmt.Sprintf("%s:%s/%s", host, port, dbName)
	}

	return fmt.Sprintf("%s/%s", host, dbName)
}
