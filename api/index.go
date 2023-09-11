package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

var srv http.Handler

func init() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	r := gin.Default()
	r.Any("/*any", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	srv = r
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
