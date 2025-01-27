package wsserver

import( 
	"net/http"
	"github.com/gorilla/websocket"
)

type Wsserver interface {
	Start() error
}

type wsSrv struct {
	mux *http.ServeMux
	srv *http.Server
	wsUpg websocket.Upgrader
}

func NewWsServer(addr string) Wsserver {
	m := http.NewServeMux()
	return  &wsSrv{
		mux: m,
		srv: &http.Server{
			Addr: addr,
			Handler: m,
		},
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.HandleFunc("/test", ws.testHandler)
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}