package wsserver

import (
	"net/http"

	"log"

	"github.com/gorilla/websocket"
)

const (
	templateDir = "./web/templates/html"
)

type WSServer interface {
	Start() error
}

type wsSrv struct {
	mux   *http.ServeMux
	srv   *http.Server
	wsUpg *websocket.Upgrader
}

func NewWsServer(addr string) WSServer {
	mux := http.NewServeMux()
	return &wsSrv{
		mux: mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsUpg: &websocket.Upgrader{},
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.HandleFunc("/test", ws.testHandler)
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Cool!"))
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error with websoket connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(conn.RemoteAddr().String())
}
