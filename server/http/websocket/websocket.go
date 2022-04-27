package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func createHandler(messageExecutor func(*websocket.Conn, []byte)) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("could not upgrade connection: %s\n", err.Error())
			return
		}

		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("could not read message: %s\n", err.Error())
				break
			}

			if messageType != websocket.TextMessage {
				log.Printf("only text messages are supported: %s\n", err.Error())
				continue
			}

			messageExecutor(conn, message)
		}
	}

	return handler
}
