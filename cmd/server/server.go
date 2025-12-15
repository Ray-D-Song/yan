package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	ServerCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Port to run the server on")
}

func startServer() {
	r := gin.Default()

	startTS := time.Now().Unix()
	r.GET("/health", func(c *gin.Context) {
		endTS := time.Now().Unix()
		start := time.Unix(startTS, 0)
		end := time.Unix(endTS, 0)

		duration := end.Sub(start)
		days := duration.Hours() / 24
		c.String(http.StatusOK, fmt.Sprintf("Service has been running for %f days", days))
	})

	// Start server on configured port
	log.Printf("Starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
