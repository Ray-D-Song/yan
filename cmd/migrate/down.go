package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	steps   int
	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations",
		Long:  `Rollback database migrations by the specified number of steps.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Rolling back %d migration(s)...\n", steps)
			// TODO: Implement migration down logic
		},
	}
)

func init() {
	downCmd.Flags().IntVarP(&steps, "steps", "s", 1, "Number of migrations to rollback")
}
