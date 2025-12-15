package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
