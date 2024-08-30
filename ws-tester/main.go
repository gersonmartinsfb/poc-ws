package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Payload struct {
	OMSId        int    `json:"OMSId,omitempty"`
	InstrumentId int    `json:"InstrumentId,omitempty"`
	Depth        int    `json:"Depth,omitempty"`
	MarketId     string `json:"MarketId,omitempty"`
}

type MessageFrame struct {
	ContentType string `json:"Content-Type"`
	UserAgent   string `json:"User-Agent"`
	M           string `json:"m"`
	I           int64  `json:"i"`
	N           string `json:"n"`
	O           string `json:"o"`
}

func connectAndSubscribe(url string, payload *Payload, wg *sync.WaitGroup) {
    defer wg.Done()

	req, _ := http.NewRequest("GET", "wss://api.foxbit.com.br/", nil)

	// Set the desired headers
	req.Header.Set("User-Agent", "ws-proxy-client/1.0")
	req.Header.Set("Content-Type", "application/json")

    // Connect to the WebSocket server
    conn, _, err := websocket.DefaultDialer.Dial(url, req.Header)
    if err != nil {
        log.Printf("Failed to connect: %v", err)
        return
    }
    defer conn.Close()

    // Subscribe to the channel
    subscribeMessage := map[string]string{
        "action":  "subscribe",
        "channel": channel,
    }
	
    if err := conn.WriteJSON(subscribeMessage); err != nil {
        log.Printf("Failed to subscribe: %v", err)
        return
    }

    // Read messages (optional, for demonstration)
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Read error: %v", err)
            break
        }
        log.Printf("Received: %s", message)
    }
}

func main() {
    url := "wss://api.foxbit.com.br/"
    channel := "test-channel"
    numConnections := 1

	payload, _ := json.Marshal(
		Payload{
			OMSId:        1,
			InstrumentId: 1,
		})

	frame := MessageFrame{
		ContentType: "application/json",
		UserAgent:   "ws-proxy-client/1.0",
		M:           "2",
		I:           1,
		N:           "SubscribeLevel2",
		O:           string(payload),
	}

    var wg sync.WaitGroup

    for i := 0; i < numConnections; i++ {
        wg.Add(1)
        go connectAndSubscribe(url, channel, &wg)
    }

    wg.Wait()
}

func (p Payload) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	// Escape double quotes in the JSON string
	data = bytes.ReplaceAll(data, []byte(`\"`), []byte(`\\\"`))

	return data, nil
}