package dbcmd

import (
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migrate-down [N]",
	Short: "Roll back N (default 1) migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateDown,
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	_ = godotenv.Load()
	cfg := config.Load()

	m, err := migrate.New("file://db/migrations", cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	n := 1
	if len(args) == 1 {
		var parseErr error
		n, parseErr = strconv.Atoi(args[0])
		if parseErr != nil {
			return fmt.Errorf("invalid N: %w", parseErr)
		}
	}

	if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate down %d: %w", n, err)
	}

	fmt.Printf("Rolled back %d migration(s).\n", n)
	return nil
}
