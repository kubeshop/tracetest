package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func createHandler(messageExecutor func(*websocket.Conn, []byte)) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		keepAlive(conn, 10*time.Second)

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

func keepAlive(conn *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	ticker := time.NewTicker(timeout / 2)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
				if err != nil {
					return
				}
				if time.Since(lastResponse) > timeout {
					conn.Close()
					return
				}
			}
		}
	}()
}
