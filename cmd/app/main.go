package main

import (
	"erpv1/config"
	"erpv1/internal/server"
	"fmt"
	"log"
)

func main() {
	run()
}

func run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error occured while parsing config file: %v", err)
	}

	server := server.NewServer(fmt.Sprintf(":%s", cfg.HttpAddr), cfg)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
