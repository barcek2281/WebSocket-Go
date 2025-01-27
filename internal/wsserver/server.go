package wsserver

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	templateDir = "./web/templates/html"
	staticDir = "./web/static/"
)

type Wsserver interface {
	Start() error
}

type wsSrv struct {
	mux   *http.ServeMux
	srv   *http.Server
	wsUpg *websocket.Upgrader
}

func NewWsServer(addr string) Wsserver {
	m := http.NewServeMux()
	return &wsSrv{
		mux: m,
		srv: &http.Server{
			Addr:    addr,
			Handler: m,
		},
		wsUpg: &websocket.Upgrader{},
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir(staticDir)) ))
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	ws.mux.HandleFunc("/test", ws.testHandler)
	return ws.srv.ListenAndServe()
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
	log.Infof("%v", conn.RemoteAddr().String())
}
