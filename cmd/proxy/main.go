package main

import (
	"GoProxy/pkg/proxy"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"strings"
)

type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Route struct {
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type Config struct {
	TLS    TLSConfig `yaml:"tls"`
	Routes []Route   `yaml:"routes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getInput(prompt, exampleValue string) string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for {
		fmt.Printf("%s (e.g. %s) : ", prompt, exampleValue)
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
	useConfig := getInput("Use config.yml? (Y/N)", "N")

	var destination, certPath, keyPath string

	if strings.ToLower(useConfig) == "y" {
		config, err := LoadConfig("config.yml")
		if err != nil {
			log.Fatal("Failed to load config.yml:", err)
		}
		destination = config.Routes[0].Target
		certPath = config.TLS.Cert
		keyPath = config.TLS.Key
	} else {
		host := getInput("Enter destination server host", "localhost")
		port := getInput("Enter destination server port", "8080")
		destination = fmt.Sprintf("http://%s:%s", host, port)

		certPath = getInput("Enter TLS certificate file path", "cert.pem")
		keyPath = getInput("Enter TLS key file path", "key.pem")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			log.Printf("[INFO] %s %s from %s (Redirect)", r.Method, r.URL.String(), r.RemoteAddr)
			return
		}
		log.Printf("[INFO] %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)
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
