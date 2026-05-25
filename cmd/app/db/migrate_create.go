package dbcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var migrationNameSanitizer = regexp.MustCompile(`[^a-z0-9_]+`)

var migrateCreateCmd = &cobra.Command{
	Use:   "migrate-create <NAME>",
	Short: "Create new up/down migration SQL files",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateCreate,
}

func runMigrateCreate(cmd *cobra.Command, args []string) error {
	name := sanitizeMigrationName(args[0])
	if name == "" {
		return fmt.Errorf("invalid NAME: must contain letters or numbers")
	}

	version, err := nextMigrationVersion(migrationDir)
	if err != nil {
		return err
	}

	baseName := fmt.Sprintf("%06d_%s", version, name)
	upPath := filepath.Join(migrationDir, baseName+".up.sql")
	downPath := filepath.Join(migrationDir, baseName+".down.sql")

	if err := os.WriteFile(upPath, []byte("-- Write your UP migration here.\n"), 0o644); err != nil {
		return fmt.Errorf("create migration file %q: %w", upPath, err)
	}
	if err := os.WriteFile(downPath, []byte("-- Write your DOWN migration here.\n"), 0o644); err != nil {
		return fmt.Errorf("create migration file %q: %w", downPath, err)
	}

	fmt.Printf("Created migration files:\n - %s\n - %s\n", upPath, downPath)
	return nil
}

func sanitizeMigrationName(raw string) string {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	normalized = migrationNameSanitizer.ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")
	for strings.Contains(normalized, "__") {
		normalized = strings.ReplaceAll(normalized, "__", "_")
	}
	return normalized
}

func nextMigrationVersion(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read migration dir: %w", err)
	}

	versions := make([]int, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}

	if len(versions) == 0 {
		return 1, nil
	}

	sort.Ints(versions)
	return versions[len(versions)-1] + 1, nil
}
