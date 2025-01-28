package main

import (
	"github.com/barcek2281/WebSocket-Go/internal/wsserver"
	log "github.com/sirupsen/logrus"
)

var addr = "192.168.16.101:8080"

func main() {
	wsSrv := wsserver.NewWsServer(addr)
	log.Info(
		"Server running: http://", addr,
	)

	err := wsSrv.Start()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
