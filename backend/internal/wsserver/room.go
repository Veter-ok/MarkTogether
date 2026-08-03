package wsserver

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	id        string
	clients   map[*websocket.Conn]struct{}
	mutex     *sync.RWMutex
	broadcast chan *wsMessage
}

func NewRoom(id string) *Room {
	return &Room{
		id:        id,
		clients:   make(map[*websocket.Conn]struct{}),
		mutex:     &sync.RWMutex{},
		broadcast: make(chan *wsMessage),
	}
}

func (room *Room) AddClient(conn *websocket.Conn) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	room.clients[conn] = struct{}{}
}

func (room *Room) RemoveClient(conn *websocket.Conn) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	delete(room.clients, conn)
}

func (room *Room) DeleteRoom() error {
	room.mutex.Lock()
	for client := range room.clients {
		if err := client.Close(); err != nil {
			return err
		}
		delete(room.clients, client)
	}
	room.mutex.Unlock()
	return nil
}

func (room *Room) Run() {
	for message := range room.broadcast {
		room.mutex.Lock()
		for client := range room.clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error with writing message: %v, in room %s", err, room.id)
			}
		}
		room.mutex.Unlock()
	}
}
