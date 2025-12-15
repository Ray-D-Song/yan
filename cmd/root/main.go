package main

import (
	"fmt"
	"os"

	"github.com/ray-d-song/yan/cmd/migrate"
	"github.com/ray-d-song/yan/cmd/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yan",
	Short: "Yan application server and migration tool",
}

func init() {
	rootCmd.AddCommand(server.ServerCmd)
	rootCmd.AddCommand(migrate.MigrateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
