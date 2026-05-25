package dbcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
)

var migrationFilenamePattern = regexp.MustCompile(`^(\d+)_.+\.(up|down)\.sql$`)

type migrationFiles struct {
	up   string
	down string
}

func currentVersion(m *migrate.Migrate) (uint, bool, error) {
	version, dirty, err := m.Version()
	if err == nil {
		return version, dirty, nil
	}
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return 0, false, err
}

func loadMigrationFiles() (map[uint]migrationFiles, error) {
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, fmt.Errorf("read migration dir: %w", err)
	}

	files := make(map[uint]migrationFiles)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		match := migrationFilenamePattern.FindStringSubmatch(name)
		if len(match) != 3 {
			continue
		}

		version, err := strconv.ParseUint(match[1], 10, 64)
		if err != nil {
			continue
		}

		fullPath := filepath.Join(migrationDir, name)
		info := files[uint(version)]
		if match[2] == "up" {
			info.up = fullPath
		} else {
			info.down = fullPath
		}
		files[uint(version)] = info
	}

	return files, nil
}

func versionsBetween(files map[uint]migrationFiles, from, to uint, direction string) []string {
	versions := make([]uint, 0, len(files))
	for version := range files {
		versions = append(versions, version)
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i] < versions[j] })

	executed := make([]string, 0)
	switch direction {
	case "up":
		for _, version := range versions {
			if version <= from || version > to {
				continue
			}
			if file := files[version].up; file != "" {
				executed = append(executed, file)
			}
		}
	case "down":
		for i := len(versions) - 1; i >= 0; i-- {
			version := versions[i]
			if version > from || version <= to {
				continue
			}
			if file := files[version].down; file != "" {
				executed = append(executed, file)
			}
		}
	}

	return executed
}

func printMigrationSummary(action string, beforeVersion, afterVersion uint, executed []string) {
	fmt.Printf("Migration action: %s (from version %d to %d)\n", action, beforeVersion, afterVersion)
	if len(executed) == 0 {
		fmt.Println("No migration files executed.")
		return
	}

	fmt.Println("Executed migration files:")
	for _, file := range executed {
		fmt.Printf(" - %s\n", file)
	}
}
