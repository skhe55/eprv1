package main

import (
	"erpv1/config"
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
}
