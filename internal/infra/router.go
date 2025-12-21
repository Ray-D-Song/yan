// Package infra provide global toolkit
package infra

import (
	"github.com/gin-gonic/gin"
)

func NewGin(config *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()

	// // Register health check endpoint under /api
	// startTS := time.Now().Unix()
	// r.GET("/api/health", func(c *gin.Context) {
	// 	endTS := time.Now().Unix()
	// 	start := time.Unix(startTS, 0)
	// 	end := time.Unix(endTS, 0)
	//
	// 	duration := end.Sub(start)
	// 	days := duration.Hours() / 24
	// 	c.String(http.StatusOK, fmt.Sprintf("Service has been running for %f days", days))
	// })

	return e
}

func NewAPIV1Group(e *gin.Engine) *gin.RouterGroup {
	v1 := e.Group("/api/v1")
	return v1
}

// // RegisterStaticFiles registers static file server using NoRoute (must be called after all API routes are registered)
// func RegisterStaticFiles(engine *gin.Engine) {
// 	sub, err := fs.Sub(embedfs.WebFile, "public")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fileServer := http.FileServer(http.FS(sub))
//
// 	// Use NoRoute to handle all unmatched routes (for SPA)
// 	engine.NoRoute(func(c *gin.Context) {
// 		// Try to serve the requested file
// 		fileServer.ServeHTTP(c.Writer, c.Request)
// 	})
// }
