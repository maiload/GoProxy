package proxy

import (
	"GoProxy/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	for _, route := range cfg.Routes {
		if strings.HasPrefix(r.URL.Path, route.Path) {
			proxyPath := strings.TrimPrefix(r.URL.Path, route.Path)
			if route.Path == "/" { proxyPath = r.URL.Path }
			proxyUrl, err := url.Parse(route.Target)
			if err != nil {
				log.Printf("[ERROR] Failed to parse target URL %q: %v", route.Target, err)
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
				return
			}

			proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				req.URL.Path = proxyPath
				req.URL.RawPath = proxyPath
				req.Host = proxyUrl.Host
			}

			proxy.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}
