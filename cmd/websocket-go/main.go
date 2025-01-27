package main

import (
	"github.com/barcek2281/WebSocket-Go/internal/wsserver"
	log "github.com/sirupsen/logrus"
)

var addr = "localhost:8080"

func main() {
	swSrv := wsserver.NewWsServer(addr)
	log.Info(
		"Server running: http://", addr,
	)
	err := swSrv.Start()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
