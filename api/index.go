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
	// 修改代理响应
	proxy.ModifyResponse = func(res *http.Response) error {
		// 添加CORS响应头
		res.Header.Set("Access-Control-Allow-Origin", "*")
		res.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		res.Header.Set("Content-Type", "text/event-stream")
		res.Header.Set("Cache-Control", "no-cache")
		res.Header.Set("Connection", "keep-alive")
		return nil
	}

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	srv = e
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
