package main

var add = make(chan TX)
var get = make(chan []TX)

type TX struct {
	TXHash        string `json:"tx_hash"`
	UnixTimestamp int64  `json:"unix_timestamp"`
}

func addTX(tx TX)  { add <- tx }
func getTXs() []TX { return <-get }

func monitor() {
	var txs []TX
	for {
		select {
		case tx := <-add:
			txs = append(txs, tx)
		case get <- txs:
			txs = []TX{}
		}
	}
}

func init() {
	go monitor()
}
