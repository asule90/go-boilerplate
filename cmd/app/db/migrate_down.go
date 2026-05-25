package dbcmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migrate-down [N]",
	Short: "Roll back N (default 1) migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateDown,
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	fmt.Printf("Running migrate-down using source %q on %s\n", migrationDir, databaseTarget(cfg.Database.URL))

	var zlog *zap.Logger
	if l, err := logger.InitLogger(cfg); err == nil {
		zlog = l
	}

	m, err := newMigrate("file://"+migrationDir, cfg.Database.URL, zlog)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Properly handle Close() errors
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close migration source: %v\n", sourceErr)
		}
		if dbErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close database connection: %v\n", dbErr)
		}
	}()

	beforeVersion, _, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version: %w", err)
	}

	files, err := loadMigrationFiles()
	if err != nil {
		return err
	}

	n := 1
	if len(args) == 1 {
		var parseErr error
		n, parseErr = strconv.Atoi(args[0])
		if parseErr != nil {
			return fmt.Errorf("invalid N: %w", parseErr)
		}
	}

	if err := m.Steps(-n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to rollback.")
			return nil
		}
		var dirtyErr migrate.ErrDirty
		if errors.As(err, &dirtyErr) {
			return fmt.Errorf("database is in a dirty state at version %d (migration failed mid-execution). Use 'app db migrate-force %d' to reset or run recovery", dirtyErr.Version, dirtyErr.Version)
		}
		if errors.Is(err, migrate.ErrLockTimeout) {
			return fmt.Errorf("could not acquire migration lock (another migration in progress?)")
		}
		return fmt.Errorf("migrate down %d: %w", n, err)
	}

	afterVersion, afterDirty, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version after down: %w", err)
	}

	if afterDirty {
		fmt.Printf("⚠️  WARNING: Database is in dirty state (version %d). Latest migration may have failed.\n", afterVersion)
	}

	executed := versionsBetween(files, beforeVersion, afterVersion, "down")
	printMigrationSummary("down", beforeVersion, afterVersion, executed)
	return nil
}
