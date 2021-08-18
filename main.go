package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const reqMessage = "{\"id\": 2, \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"new_transaction\"]}"
const nodeID = 1 //TODO: id and coordinate

var firstRecv = true
var data []TX

type Recv struct {
	Jsonrpc string
	Method  string
	Params  struct {
		Result       string
		Subscription string
	}
}

type TX struct {
	NodeID        int    `json:"node_id"`
	TXHash        string `json:"tx_hash"`
	UnixTimestamp int64  `json:"unix_timestamp"`
}

func sendData() {

	ticker := time.NewTicker(time.Minute).C

	for {
		select {
		case <-ticker:
			if len(data) == 0 {
				continue
			}
			jsonValue, _ := json.Marshal(data)

			response, err := http.Post("http://localhost:3006/txs", "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				log.Fatal("post:", err)
			}

			response.Body.Close()
			data = []TX{}
		}
	}
}

func getHash(message []byte) string {
	var recv Recv
	if err := json.Unmarshal(message, &recv); err != nil {
		log.Fatal("unmarshal into Recv:", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(recv.Params.Result), &jsonMap); err != nil {
		log.Fatal("unmarshal into jsonMap:", err)
	}

	transaction := jsonMap["transaction"]
	transactionMap := transaction.(map[string]interface{})

	return transactionMap["hash"].(string)
}

func main() {
	go sendData()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:18115", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(reqMessage)); err != nil {
		log.Fatal("write:", err)
	}

	for {
		if _, message, err := conn.ReadMessage(); err != nil {
			log.Fatal("read:", err)
		} else {
			if firstRecv {
				firstRecv = false
				continue
			}

			tx := TX{nodeID, getHash(message), time.Now().Unix()}
			fmt.Println(tx)

			data = append(data, tx)
		}
	}
}
