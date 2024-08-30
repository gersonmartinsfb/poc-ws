package main

import (
    "log"
    "sync"

    "github.com/gorilla/websocket"
)

func connectAndSubscribe(url string, channel string, wg *sync.WaitGroup) {
    defer wg.Done()

    // Connect to the WebSocket server
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
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
    url := "ws://localhost:8080/ws"
    channel := "test-channel"
    numConnections := 10

    var wg sync.WaitGroup

    for i := 0; i < numConnections; i++ {
        wg.Add(1)
        go connectAndSubscribe(url, channel, &wg)
    }

    wg.Wait()
}