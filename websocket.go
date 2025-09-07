package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WSServer struct {
	Log        *zap.Logger
	Clients    map[*websocket.Conn]bool
	ClientsMux sync.Mutex
}

func NewWSServer() *WSServer {
	ws := WSServer{
		Clients: map[*websocket.Conn]bool{},
	}
	return &ws
}

func (ws *WSServer) Handler() http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ws.Log.Error("Failed to upgrade connection", zap.Error(err))
			return
		}
		defer conn.Close()

		ws.ClientsMux.Lock()
		ws.Clients[conn] = true
		ws.ClientsMux.Unlock()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				ws.Log.Info("Client disconnected", zap.Error(err))
				break
			}
		}

		ws.ClientsMux.Lock()
		delete(ws.Clients, conn)
		ws.ClientsMux.Unlock()
	}
}

func (ws *WSServer) Notify(path string) {
	ws.Log.Debug("Notifying clients", zap.String("path", path))
	ws.ClientsMux.Lock()
	defer ws.ClientsMux.Unlock()

	for conn := range ws.Clients {
		msg := map[string]string{"path": path}
		if err := conn.WriteJSON(msg); err != nil {
			ws.Log.Error("Failed to send message to client", zap.Error(err))
			delete(ws.Clients, conn)
		}
	}
}
