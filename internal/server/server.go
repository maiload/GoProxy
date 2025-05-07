package server

import (
	"GoProxy/config"
	"GoProxy/internal/proxy"
	"log"
	"net/http"
)

func Start(cfg *config.Config) {
	sslPort, certPath, keyPath := cfg.Server.SSL.Port, cfg.Server.SSL.Cert, cfg.Server.SSL.Key
	isHTTPS := sslPort != "" && certPath != "" && keyPath != ""

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil && isHTTPS {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusPermanentRedirect)
			log.Printf("[INFO] %s %s from %s (Redirect)", r.Method, r.URL.String(), r.RemoteAddr)
			return
		}
		log.Printf("[INFO] %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)
		proxy.Handle(w, r, cfg)
	})

	if isHTTPS {
		go func() {
			log.Println("Starting HTTPS server on :" + sslPort)
			if err := http.ListenAndServeTLS(":"+sslPort, certPath, keyPath, nil); err != nil {
				log.Fatal("HTTPS Server Error:", err)
			}
		}()
	} else {
		log.Println("The values for certPath and keyPath are empty.")
	}

	port := cfg.Server.Port
	log.Println("Starting HTTP server on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("HTTP Server Error:", err)
	}
}
