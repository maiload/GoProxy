package server

import (
	"GoProxy/config"
	"GoProxy/internal/proxy"
	"log"
	"net/http"
	"sync"
)

func Start(cfg *config.Config) {
	port, sslPort, certPath, keyPath := cfg.Server.Port, cfg.Server.SSL.Port, cfg.Server.SSL.Cert, cfg.Server.SSL.Key
	isHTTPS := sslPort != "" && certPath != "" && keyPath != ""
	isHTTP := port != ""

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil && isHTTPS {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusPermanentRedirect)
			log.Printf("[INFO] %s %s from %s (Redirect)", r.Method, r.URL.String(), r.RemoteAddr)
			return
		}
		log.Printf("[INFO] %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)
		proxy.Handle(w, r, cfg)
	})

	var wg sync.WaitGroup
	if isHTTPS {
		wg.Add(1)
		go func() {
			log.Println("Starting HTTPS server on :" + sslPort)
			if err := http.ListenAndServeTLS(":"+sslPort, certPath, keyPath, nil); err != nil {
				log.Fatal("HTTPS Server Error:", err)
			}
		}()
	} else {
		log.Println("The values for ssl are empty.")
	}

	if isHTTP {
		log.Println("Starting HTTP server on :" + port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal("HTTP Server Error:", err)
		}
	} else {
		log.Println("The value of the port is empty.")
		wg.Wait()
	}
}
