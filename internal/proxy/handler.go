package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func IndexPathHandler(w http.ResponseWriter, r *http.Request, dest string) {
	proxyUrl, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(w, r)
}
