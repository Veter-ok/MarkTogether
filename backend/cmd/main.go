package main

import (
	"log"

	"github.com/Veter-ok/MarkTogether/internal/app"
)

const (
	addr = "localhost:8080"
)

func main() {
	app := app.NewApp(addr)
	if err := app.Start(); err != nil {
		log.Printf("Error with server starting: %v", err)
	}
}
