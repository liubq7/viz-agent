package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const reqNodeInfo = "{\"id\": 42, \"jsonrpc\": \"2.0\", \"method\": \"local_node_info\", \"params\": []}"
const reqNewTX = "{\"id\": 2, \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"new_transaction\"]}"

var firstRecv = true
var id = make(chan string)

type Recv struct {
	Jsonrpc string
	Method  string
	Params  struct {
		Result       string
		Subscription string
	}
}

func setNodeID() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:18115", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(reqNodeInfo)); err != nil {
		log.Fatal("req node info:", err)
	}

	if _, message, err := conn.ReadMessage(); err != nil {
		log.Fatal("read:", err)
	} else {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(message, &jsonMap); err != nil {
			log.Fatal("unmarshal into jsonMap:", err)
		}

		result := jsonMap["result"]
		resultMap := result.(map[string]interface{})
		nodeID := resultMap["node_id"].(string)

		id <- nodeID
		fmt.Println(nodeID)
		sendNodeID(nodeID)
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
	go setNodeID()
	go sendTXs()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://localhost:18115", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer conn.Close()

	if err := conn.WriteMessage(websocket.TextMessage, []byte(reqNewTX)); err != nil {
		log.Fatal("req new tx:", err)
	}

	nodeID := <- id

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
			addTX(tx)
		}
	}
}
