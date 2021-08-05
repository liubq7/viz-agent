package main

import (
	"log"

	"github.com/gorilla/websocket"
)

const reqMessage = "{\"id\": 2, \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"new_transaction\"]}"

func main() {
	dialer := websocket.Dialer{}

	conn, _, err := dialer.Dial("ws://localhost:18115", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(reqMessage)); err != nil {
		log.Fatal(err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
		log.Printf("recv: %s", message)
	}
}
