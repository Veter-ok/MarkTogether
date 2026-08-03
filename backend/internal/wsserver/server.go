package wsserver

import (
	"net"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/Veter-ok/MarkTogether/internal/document"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	wsUpg   *websocket.Upgrader
	wsRooms map[string]*Room
	mutex   *sync.RWMutex
	store   document.Store
}

func NewWsServer(addr string, store document.Store) *WSServer {
	return &WSServer{
		wsUpg:   &websocket.Upgrader{},
		wsRooms: make(map[string]*Room),
		mutex:   &sync.RWMutex{},
		store:   store,
	}
}

func (ws *WSServer) Stop() {
	ws.mutex.Lock()
	for id, room := range ws.wsRooms {
		if err := room.DeleteRoom(); err != nil {
			log.Printf("Error with closing room %s: %v", id, err)
		}
		delete(ws.wsRooms, id)
	}
	ws.mutex.Unlock()
}

func (ws *WSServer) getOrCreateRoom(roomID string) (*Room, error) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if room, ok := ws.wsRooms[roomID]; ok {
		return room, nil
	}

	doc, err := ws.store.Get(roomID)
	if err != nil || doc == nil {
		return nil, err
	}

	room := NewRoom(roomID)
	go room.Run()
	ws.wsRooms[roomID] = room
	log.Printf("Room %s created", roomID)
	return room, nil
}

func (ws *WSServer) WSHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "missing room parameter", http.StatusBadRequest)
		return
	}

	room, err := ws.getOrCreateRoom(roomID)
	if err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}

	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	log.Printf("Client %s connected to room %s", conn.RemoteAddr().String(), roomID)
	room.AddClient(conn)
	go ws.readFromClient(conn, room)
}

func (ws *WSServer) readFromClient(conn *websocket.Conn, room *Room) {
	for {
		msg := new(wsMessage)
		if err := conn.ReadJSON(msg); err != nil {
			wsErr, ok := err.(*websocket.CloseError)
			if !ok || wsErr.Code != websocket.CloseGoingAway {
				log.Printf("Error with reading from WebSoket: %v", err)
			}
			break
		}
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Printf("Error with address split: %v", err)
			break
		}
		msg.IPAdress = host
		msg.Time = time.Now().Format("15:04")
		room.broadcast <- msg
	}
	room.RemoveClient(conn)
	conn.Close()
	log.Printf("Client %s disconnected from room %s", conn.RemoteAddr().String(), room.id)
}
