package dbcmd

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

var migrateForceCmd = &cobra.Command{
	Use:   "migrate-force <VERSION>",
	Short: "Force migration version without running SQL",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateForce,
}

func runMigrateForce(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	fmt.Printf("Running migrate-force using source %q on %s\n", migrationDir, databaseTarget(cfg.Database.URL))

	version, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid VERSION: %w", err)
	}

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

	if err := m.Force(version); err != nil {
		return fmt.Errorf("force version %d: %w", version, err)
	}

	afterVersion, afterDirty, err := currentVersion(m)
	if err != nil {
		return fmt.Errorf("read migration version after force: %w", err)
	}

	fmt.Printf("Migration state updated: version %d (dirty=%t) → %d (dirty=%t)\n", beforeVersion, beforeDirty, afterVersion, afterDirty)

	if beforeDirty && !afterDirty {
		fmt.Printf("✓ Recovered from dirty state\n")
	}

	fmt.Println("Note: Force only updates schema history state without executing migration files.")
	return nil
}
