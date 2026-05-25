package dbcmd

import "github.com/spf13/cobra"

// DBCmd is the parent command for database operations.
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
}

func init() {
	DBCmd.AddCommand(migrateUpCmd)
	DBCmd.AddCommand(migrateDownCmd)
	DBCmd.AddCommand(migrateCreateCmd)
	DBCmd.AddCommand(migrateForceCmd)
}
