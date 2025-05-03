package main

import (
	"GoProxy/pkg/proxy"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getInput(prompt, exampleValue string) string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for {
		fmt.Printf("%s (e.g., %s) : ", prompt, exampleValue)
		scanner.Scan()
		input = scanner.Text()
		if input != "" {
			break
		}
		fmt.Println("Input cannot be empty. Please try again.")
	}
	return input
}

func main() {
	host := getInput("Enter destination server host", "localhost")
	port := getInput("Enter destination server port", "8080")
	destination := fmt.Sprintf("http://%s:%s", host, port)

	certPath := getInput("Enter TLS certificate file path", "cert.pem")
	keyPath := getInput("Enter TLS key file path", "key.pem")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		}
		proxy.IndexPathHandler(w, r, destination)
	})

	go func() {
		log.Println("Starting HTTP server on :80 to redirect to HTTPS...")
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatal("HTTP Server Error:", err)
		}
	}()

	log.Println("Starting proxy server on :443")
	if err := http.ListenAndServeTLS(":443", certPath, keyPath, nil); err != nil {
		log.Fatal("HTTPS Server Error:", err)
	}
}
