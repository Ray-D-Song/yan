// Package server provides the web server command for the Yan application.
package server

import (
	"github.com/ray-d-song/yan/internal/app"
	"github.com/spf13/cobra"
)

var port string

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", ":18080", "Port to run the server on")
}

func startServer() {
	app.New().Run()
}
