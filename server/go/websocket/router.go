package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type MessageExecutor func(*websocket.Conn, Message)

type Router struct {
	routes map[string]MessageExecutor
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]MessageExecutor),
	}
}

func (r *Router) Add(messageType string, executor MessageExecutor) {
	r.routes[messageType] = executor
}

func (r *Router) ListenAndServe(addr string) {
	routingFunction := func(conn *websocket.Conn, message []byte) {
		messageObject := Message{}
		err := json.Unmarshal(message, &messageObject)
		if err != nil {
			errMessage := Error(fmt.Errorf("could not unmarshal message: %w", err))
			conn.WriteJSON(errMessage)
			return
		}

		if handler, exists := r.routes[messageObject.Type]; exists {
			handler(conn, messageObject)
		} else {
			conn.WriteJSON(Error(fmt.Errorf("No routes for message type %s", messageObject.Type)))
		}
	}

	http.HandleFunc("/ws", createHandler(routingFunction))
	log.Fatal(http.ListenAndServe(addr, nil))
}
