package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Payload struct {
	OMSId        int    `json:"OMSId,omitempty"`
	InstrumentId int    `json:"InstrumentId,omitempty"`
	Depth        int    `json:"Depth,omitempty"`
	MarketId     string `json:"MarketId,omitempty"`
}

type MessageFrame struct {
	M int64  `json:"m"`
	I int64  `json:"i"`
	N string `json:"n"`
	O string `json:"o"`
}

func connectAndSubscribe(url string, payload *Payload, wg *sync.WaitGroup) {
	defer wg.Done()

	req, _ := http.NewRequest("GET", "wss://api.foxbit.com.br/", nil)

	// Set the desired headers
	req.Header.Set("User-Agent", "ws-tester/1.0")
	req.Header.Set("Content-Type", "application/json")

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	frame := MessageFrame{
		M: 2,
		I: time.Now().UnixNano(),
		N: "SubscribeLevel2",
		O: string(payloadJson),
	}

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(url, req.Header)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer conn.Close()

	// Subscribe to the channel
	time.Sleep(1 * time.Second)
	if err := conn.WriteJSON(frame); err != nil {
		log.Printf("Failed to subscribe: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("(%s) Received: %s", payload.MarketId, message)
	}
}

func main() {
	url := "wss://api.foxbit.com.br/"

	var wg sync.WaitGroup

	markets := []string{"btcbrl", "ethbrl", "xrpbrl", "adabrl", "apebrl", "nexobrl", "sushibrl", "1inchbrl", "atombrl", "trxbrl", "tonbrl"}
	for _, market := range markets {
		wg.Add(1)
		go connectAndSubscribe(url, &Payload{
			OMSId:    1,
			MarketId: market,
			Depth:    10,
		}, &wg)
	}

	wg.Wait()
}
