package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// var srv *gin.Engine // 注意这里改为 *gin.Engine

// func init() {
// 	proxyURL := os.Getenv("Proxy_URL")
// 	target, _ := url.Parse(proxyURL)
// 	proxy := httputil.NewSingleHostReverseProxy(target)

// 	r := gin.Default()
// 	r.Any("/*any", func(c *gin.Context) {
// 		proxy.ServeHTTP(c.Writer, c.Request)
// 	})
// 	srv = r
// }

// func Proxy(w http.ResponseWriter, r *http.Request) {
// 	srv.ServeHTTP(w, r)
// }

var (
	app *gin.Engine
)

// CREATE ENDPOIND

func myRoute(r *gin.RouterGroup) {
	// r.GET("/admin", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "Hello from golang in vercel")
	// })

	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	r.Any("/*any", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func init() {
	app = gin.New()
	// r := app.Group("/api")
	r := app.Group("/")
	myRoute(r)

}

// Handler ADD THIS SCRIPT
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
