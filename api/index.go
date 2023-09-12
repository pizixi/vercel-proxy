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
	// proxyURL := os.Getenv("Proxy_URL")
	// target, _ := url.Parse(proxyURL)
	// proxy := httputil.NewSingleHostReverseProxy(target)

	// e := echo.New()
	// e.Any("/*", func(c echo.Context) error {
	// 	proxy.ServeHTTP(c.Response().Writer, c.Request())
	// 	return nil
	// })
	e := echo.New()

	proxySites := map[string]string{
		"/reproxy-workergpt/*":    "https://workergpt.cn",
		"/reproxy-eqing/*":        "https://next.eqing.tech",
		"/reproxy-chatanywhere/*": "https://api.chatanywhere.cn",
		"/reproxy-lbbai/*":        "https://postapi.lbbai.cc",
		"/reproxy-mixerbox/*":     "https://chatai.mixerbox.com",
	}

	modifyProxyResponse := func(res *http.Response) error {
		res.Header.Set("Access-Control-Allow-Origin", "*")
		res.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		res.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		res.Header.Set("Content-Type", "text/event-stream")
		res.Header.Set("Cache-Control", "no-cache")
		res.Header.Set("Connection", "keep-alive")
		return nil
	}

	for routePath, targetURL := range proxySites {
		target, _ := url.Parse(targetURL)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ModifyResponse = modifyProxyResponse

		e.Any(routePath, func(c echo.Context) error {
			c.Request().URL.Path = target.Path + c.Param("*")
			c.Request().Host = target.Host
			proxy.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
	}

	e.Any("/*", func(c echo.Context) error {
		target, _ := url.Parse(os.Getenv("Proxy_URL"))
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	srv = e
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	srv.ServeHTTP(w, r)
}
