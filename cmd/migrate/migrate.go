// Package migrate provides database migration commands for the Yan application.
package migrate

import (
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage database migrations with up, down, and new commands.`,
}

func init() {
	MigrateCmd.AddCommand(upCmd)
	MigrateCmd.AddCommand(downCmd)
	MigrateCmd.AddCommand(newCmd)
}
