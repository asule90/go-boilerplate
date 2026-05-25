package dbcmd

import (
	"fmt"
	"net/url"
	"strings"
)

const migrationDir = "db/migrations"

func databaseTarget(databaseURL string) string {
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
