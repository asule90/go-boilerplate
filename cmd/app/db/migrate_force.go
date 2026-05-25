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

	m, err := migrate.New("file://"+migrationDir, cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

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

	fmt.Printf("Migration action: force (from version %d dirty=%t to version %d dirty=%t)\n", beforeVersion, beforeDirty, afterVersion, afterDirty)
	fmt.Println("No migration files executed. Force only updates schema history state.")
	return nil
}
