package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/poller"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allowing all origins for dev/testing
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients map[*Client]bool
	mu      sync.RWMutex
	poller  *poller.Poller
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
}

func NewHub(p *poller.Poller) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		poller:  p,
	}
}

func (h *Hub) broadcastTrains() {
	trains := h.poller.GetTrains()
	data, err := json.Marshal(trains)
	if err != nil {
		log.Printf("Error marshaling trains: %v", err)
		return
	}

	h.mu.RLock()
	for client := range h.clients {
		if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Error writing to client: %v", err)
			h.mu.RUnlock()
			h.removeClient(client)
			h.mu.RLock()
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		client.conn.Close()
		delete(h.clients, client)
	}
}

func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	h.broadcastTrains()
}
