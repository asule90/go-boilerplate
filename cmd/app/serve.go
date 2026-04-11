package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/internal/db"
	"github.com/sule/go-boilerplate/internal/middleware"
	"github.com/sule/go-boilerplate/internal/server"
	"github.com/sule/go-boilerplate/pkg/logger"
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

	pool, err := db.NewPgxPool(context.Background(), cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

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
			log.Fatal("server error")
		}
	}()

	<-quit
	return app.Shutdown()
}
