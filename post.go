package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func post(url string, jsonValue []byte) {
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal("post:", err)
	}

	response.Body.Close()
}

func sendNodeID(nodeID string) {
	values := map[string]string{"node_id": nodeID}
	jsonValue, _ := json.Marshal(values)
	post("http://localhost:3006/nodes", jsonValue)
}

func sendTXs() {
	ticker := time.NewTicker(time.Minute).C

	for {
		select {
		case <-ticker:
			txs := getTXs()
			if len(txs) == 0 {
				continue
			}
			jsonValue, _ := json.Marshal(txs)
			post("http://localhost:3006/txs", jsonValue)
		}
	}
}
