package infra

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(config *Config) *gin.Engine {
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

	return r
}
