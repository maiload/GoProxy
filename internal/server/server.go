package server

import (
	"GoProxy/config"
	"GoProxy/internal/proxy"
	"log"
	"net/http"
)

func Start(cfg *config.Config) {
	certPath, keyPath := cfg.Server.SSL.Cert, cfg.Server.SSL.Key
	isHTTPS := certPath != "" && keyPath != ""

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil && isHTTPS {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			log.Printf("[INFO] %s %s from %s (Redirect)", r.Method, r.URL.String(), r.RemoteAddr)
			return
		}
		log.Printf("[INFO] %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)
		proxy.Handle(w, r, cfg)
	})

	port := cfg.Server.Port

	if isHTTPS {
		go func() {
			log.Println("Starting HTTPS server on :" + port)
			if err := http.ListenAndServeTLS(":"+port, certPath, keyPath, nil); err != nil {
				log.Fatal("HTTPS Server Error:", err)
			}
		}()
	} else {
		log.Println("The values for certPath and keyPath are empty.")
	}

	httpPort := port
	if port == "443" {
		httpPort = "80"
	}
	log.Println("Starting HTTP server on :" + httpPort)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatal("HTTP Server Error:", err)
	}
}
