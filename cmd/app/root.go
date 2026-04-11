package app

import (
	"os"

	"github.com/spf13/cobra"
	dbcmd "github.com/sule/go-boilerplate/cmd/app/db"
)

var rootCmd = &cobra.Command{
	Use:   "boilerplate",
	Short: "Go Fiber boilerplate application",
	Long:  "A production-ready Go backend boilerplate using Fiber, PostgreSQL, and more.",
}

var dbCmd = dbcmd.DBCmd

// Execute runs the root cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(dbCmd)
}
