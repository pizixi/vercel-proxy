package handler

import (
	"io"
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

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		// 设置超时时间
		timeout := time.Duration(50) * time.Second
		ctx := c.Request().Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		c.SetRequest(c.Request().WithContext(ctx))

		// 创建目标响应的主体
		pipeReader, pipeWriter := io.Pipe()
		defer pipeReader.Close()

		// 修改源和目标的响应主体
		c.Response().Writer = pipeWriter
		proxy.Transport = &http.Transport{
			Proxy: http.ProxyURL(target),
		}

		// 复制源响应的主体到目标响应的主体
		go func() {
			defer pipeWriter.Close()
			proxy.ServeHTTP(c.Response().Writer, c.Request())
		}()

		// 将目标响应的主体复制到客户端的响应主体
		io.Copy(c.Response().Writer, pipeReader)

		return nil
	})
	srv = e
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
