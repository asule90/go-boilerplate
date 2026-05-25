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

var migrateUpCmd = &cobra.Command{
	Use:   "migrate-up [N]",
	Short: "Apply N (or all) pending migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateUp,
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	fmt.Printf("Running migrate-up using source %q on %s\n", migrationDir, databaseTarget(cfg.Database.URL))

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

	beforeVersion, beforeDirty, err := currentVersion(m)
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
		if err := m.Steps(n); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to apply.")
				return nil
			}
			var dirtyErr migrate.ErrDirty
			if errors.As(err, &dirtyErr) {
				return fmt.Errorf("database is in a dirty state at version %d (migration failed mid-execution). Use 'app db migrate-force %d' to reset or run recovery", dirtyErr.Version, dirtyErr.Version)
			}
			if errors.Is(err, migrate.ErrLockTimeout) {
				return fmt.Errorf("could not acquire migration lock (another migration in progress?)")
			}
			return fmt.Errorf("migrate up %d: %w", n, err)
		}
	} else {
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to apply.")
				return nil
			}
			var dirtyErr migrate.ErrDirty
			if errors.As(err, &dirtyErr) {
				return fmt.Errorf("database is in a dirty state at version %d (migration failed mid-execution). Use 'app db migrate-force %d' to reset or run recovery", dirtyErr.Version, dirtyErr.Version)
			}
			if errors.Is(err, migrate.ErrLockTimeout) {
				return fmt.Errorf("could not acquire migration lock (another migration in progress?)")
			}
			return fmt.Errorf("migrate up: %w", err)
		}
	}

	afterVersion, afterDirty, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version after up: %w", err)
	}

	if afterDirty {
		fmt.Printf("⚠️  WARNING: Database is in dirty state (version %d). Latest migration may have failed.\n", afterVersion)
	}

	if beforeDirty && !afterDirty {
		fmt.Printf("✓ Recovered from dirty state (version %d -> %d)\n", beforeVersion, afterVersion)
	}

	executed := versionsBetween(files, beforeVersion, afterVersion, "up")
	printMigrationSummary("up", beforeVersion, afterVersion, executed)
	return nil
}
