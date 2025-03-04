package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func newWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (hub *WebSocketHub) handleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return err
	}
	defer conn.Close()

	hub.mu.Lock()
	hub.clients[conn] = true
	hub.mu.Unlock()

	log.Println("New WebSocket connection established")

	// Handle disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			hub.mu.Lock()
			delete(hub.clients, conn)
			hub.mu.Unlock()
			log.Println("WebSocket client disconnected")
			break
		}
	}
	return nil
}

func (hub *WebSocketHub) broadcastMessage(message string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for client := range hub.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending message:", err)
			client.Close()
			delete(hub.clients, client)
		}
	}
}
