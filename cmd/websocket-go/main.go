package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/barcek2281/WebSocket-Go/internal/wsserver"
)

func main() {
	addr := "localhost:8080"
	swSrv := wsserver.NewWsServer(addr)
	log.Info(
		"Server running: http://", addr,
	)
	err := swSrv.Start()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
