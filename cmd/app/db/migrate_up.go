package dbcmd

import (
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
)

var migrateUpCmd = &cobra.Command{
	Use:   "migrate-up [N]",
	Short: "Apply N (or all) pending migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateUp,
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	fmt.Printf("Running migrate-up using source %q on %s\n", migrationDir, databaseTarget(cfg.Database.URL))

	m, err := migrate.New("file://"+migrationDir, cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	beforeVersion, _, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version: %w", err)
	}

	files, err := loadMigrationFiles()
	if err != nil {
		return err
	}

	if len(args) == 1 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid N: %w", err)
		}
		if err := m.Steps(n); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrate up %d: %w", n, err)
		}
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrate up: %w", err)
		}
	}

	afterVersion, _, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version after up: %w", err)
	}

	executed := versionsBetween(files, beforeVersion, afterVersion, "up")
	printMigrationSummary("up", beforeVersion, afterVersion, executed)
	return nil
}
