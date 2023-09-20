package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

var srv http.Handler

func init() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		// 设置请求的Host头
		c.Request().Host = target.Host
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	srv = e
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
