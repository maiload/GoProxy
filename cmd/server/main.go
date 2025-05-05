package main

import (
	"GoProxy/config"
	"GoProxy/internal/server"
	"log"
)

func main() {
	log.Println("Starting proxy server...")

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal("Failed to load config.yml:", err)
	}

	server.Start(cfg)
}
