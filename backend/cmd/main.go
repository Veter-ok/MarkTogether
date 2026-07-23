package main

import (
	"log"

	"github.com/Veter-ok/MarkTogether/internal/wsserver"
)

const (
	addr = "localhost:8080"
)

func main() {
	wsSrv := wsserver.NewWsServer(addr)
	if err := wsSrv.Start(); err != nil {
		log.Printf("Error with server starting: %v", err)
	}
}
