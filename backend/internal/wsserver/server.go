package wsserver

import (
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/gorilla/websocket"
)

const (
	templateDir = "../web/templates/html"
)

type WSServer interface {
	Start() error
}

type wsSrv struct {
	mux       *http.ServeMux
	srv       *http.Server
	wsUpg     *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	mutex     *sync.RWMutex
	broadcast chan *wsMessage
}

func NewWsServer(addr string) WSServer {
	mux := http.NewServeMux()
	return &wsSrv{
		mux: mux,
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsUpg:     &websocket.Upgrader{},
		wsClients: map[*websocket.Conn]struct{}{},
		mutex:     &sync.RWMutex{},
		broadcast: make(chan *wsMessage),
	}
}

func (ws *wsSrv) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	ws.mux.HandleFunc("/test", ws.testHandler)
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	go ws.writeToClientsBroadcast()
	return ws.srv.ListenAndServe()
}

func (ws *wsSrv) testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Cool!"))
}

func (ws *wsSrv) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error with websoket connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Client with address %s connected", conn.RemoteAddr().String())
	ws.mutex.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mutex.Unlock()
	go ws.readFromClient(conn)
}

func (ws *wsSrv) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(wsMessage)
		if err := conn.ReadJSON(msg); err != nil {
			log.Printf("Error with websoket connection: %v", err)
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Printf("Error with address split: %v", err)
		}
		msg.IPAdress = host
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
					log.Printf("Error with writing message: %v", err)
				}
			}()
		}
		ws.mutex.RUnlock()
	}
}
