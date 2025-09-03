package websocket

import (
	"fmt"
	"log"

	"github.com/Rx-11/hft-arbitrage-engine/internal/handler"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn    *websocket.Conn
	apiKey  string
	url     string
	symbols []string
}

func NewWebSocketClient(apiKey string, symbols []string) *WebSocketClient {
	return &WebSocketClient{
		apiKey:  apiKey,
		url:     fmt.Sprintf("wss://ws.finnhub.io?token=%s", apiKey),
		symbols: symbols,
	}
}

func (client *WebSocketClient) Connect() error {
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	log.Println("Connected to Finnhub WebSocket")
	return nil
}

func (client *WebSocketClient) Subscribe() {
	for _, symbol := range client.symbols {
		message := fmt.Sprintf(`{"type":"subscribe","symbol":"%s"}`, symbol)
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("error subscribing to %s: %v", symbol, err)
		} else {
			log.Printf("Subscribed to %s", symbol)
		}
	}
}

func (client *WebSocketClient) Listen() {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		// log.Printf("Received message: %s", message)
		handler.ProcessMarketData(message)
	}
}

func (client *WebSocketClient) Close() {
	err := client.conn.Close()
	if err != nil {
		log.Printf("Error closing WebSocket connection: %v", err)
	}
	log.Println("WebSocket connection closed")
}
