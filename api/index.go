package handler

import (
	"net"
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

	// 设置流式输出和超时时间
	proxy.FlushInterval = time.Millisecond * 100 // 流式输出的间隔，这里设置为100毫秒，可以根据需求调整
	proxy.Transport = &http.Transport{           // 设置超时时间
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second, // 超时时间，从环境变量中读取
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
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
