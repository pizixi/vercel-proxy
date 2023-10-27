package handler

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

var srv http.Handler

func init() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 设置流式输出
	proxy.FlushInterval = time.Millisecond * 100

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		// 设置请求超时
		ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(50)*time.Second)
		defer cancel()
		c.SetRequest(c.Request().WithContext(ctx))

		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	srv = e
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
