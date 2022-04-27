package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kubeshop/tracetest/subscription"
)

type MessageExecutor interface {
	Execute(*websocket.Conn, []byte)
}

type routingMessage struct {
	Type string `json:"type"`
}

type Router struct {
	routes              map[string]MessageExecutor
	subscriptionManager *subscription.Manager
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
		messageObject := routingMessage{}
		err := json.Unmarshal(message, &messageObject)
		if err != nil {
			errMessage := ErrorMessage(fmt.Errorf("could not unmarshal message: %w", err))
			conn.WriteJSON(errMessage)
			return
		}

		if messageExecutor, exists := r.routes[messageObject.Type]; exists {
			messageExecutor.Execute(conn, message)
		} else {
			conn.WriteJSON(ErrorMessage(fmt.Errorf("No routes for message type %s", messageObject.Type)))
		}
	}

	http.HandleFunc("/ws", createHandler(routingFunction))
	log.Fatal(http.ListenAndServe(addr, nil))
}
