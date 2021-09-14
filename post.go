package main

import (
	"bytes"
	"encoding/json"
	"flag"
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

func sendTXs(id <-chan string) {
	urlPtr := flag.String("url", "http://54.254.68.135", "url")
	flag.Parse()

	ticker := time.NewTicker(time.Minute).C
	nodeID := <-id

	for {
		select {
		case <-ticker:
			txs := getTXs()
			if len(txs) == 0 {
				continue
			}
			jsonValue, _ := json.Marshal(txs)
			post(*urlPtr+"/api/nodes/"+nodeID+"/info", jsonValue)
		}
	}
}
