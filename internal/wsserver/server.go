package wsserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	templateDir = "./web/templates/html"
	staticDir   = "./web/static/"
)

type Wsserver interface {
	Start() error
	Stop() error
}

type wsSrv struct {
	mux       *http.ServeMux
	srv       *http.Server
	wsUpg     *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	mutex     *sync.RWMutex
	broadcast chan *wsMessage
}

func NewWsServer(addr string) Wsserver {
	m := http.NewServeMux()
	return &wsSrv{
		mux: m,
		srv: &http.Server{
			Addr:    addr,
			Handler: m,
		},
		wsUpg:     &websocket.Upgrader{},
		wsClients: map[*websocket.Conn]struct{}{},
		mutex:     &sync.RWMutex{},
		broadcast: make(chan *wsMessage),
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	ws.mux.HandleFunc("/test", ws.testHandler)
	go ws.writeToClientsBroadcast()
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) Stop() error {
	log.Infof("server stop working")
	close(ws.broadcast)
	ws.mutex.Lock()
	for conn := range ws.wsClients {
	 	conn.Close()
		delete(ws.wsClients, conn)
	}
	return ws.srv.Shutdown(context.Background())
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Salem alem"))
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error with websocket conn: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Infof("Clients with this addr: %v, connected", conn.RemoteAddr().String())
	ws.mutex.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mutex.Unlock()
	go ws.readFromClient(conn)
}

func (ws *wsSrv) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(wsMessage)
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Errorf("Error with reading from client ws, %v", err)
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Errorf("Erro with host and port: %v", err)
		}
		msg.IpAddress = host
		msg.Time = time.Now().Format("15:04")
		ws.broadcast <- msg
	}
	ws.mutex.Lock()
	delete(ws.wsClients, conn)
	ws.mutex.Unlock()
}

func (ws *wsSrv) writeToClientsBroadcast() {
	for msg := range ws.broadcast {
		ws.mutex.RLock()
		for client := range ws.wsClients {
			 func() {
				if err := client.WriteJSON(msg); err != nil {
					log.Errorf("Error write a message: %v", err)
				}
			}()
		}
		ws.mutex.RUnlock()
	}
}
