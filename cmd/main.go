package main

import (
	"library-api/internal/app"
	"library-api/pkg/config"
	"log"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatalln("failed to load configs:", err)
	}
}

func main() {
	app.Start()
}
