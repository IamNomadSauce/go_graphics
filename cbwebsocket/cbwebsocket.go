package cbwebsocket

import (
  "fmt"
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// WebsocketClient represents a WebSocket client
type WebsocketClient struct {
	Conn       *websocket.Conn
	ProductIDs []string
	Channels   []string
	MessageChannel chan string
}

// NewWebsocketClient creates a new WebSocket client
func NewWebsocketClient(productIDs []string, websocketURI string, channels []string, messageChannel chan string) *WebsocketClient {
	if websocketURI == "" {
		websocketURI = "wss://ws-feed.pro.coinbase.com"
	}
	if len(channels) == 0 {
		channels = []string{"full", "heartbeat", "ticker"}
	}
	return &WebsocketClient{
		ProductIDs: productIDs,
		Channels:   channels,
		MessageChannel: messageChannel,
	}
}

// Connect establishes a WebSocket connection
func (c *WebsocketClient) Connect(websocketURI string) {
	u, err := url.Parse(websocketURI)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	c.Conn = conn
}

// Subscribe sends a subscription message to the WebSocket server
func (c *WebsocketClient) Subscribe() {
	subscription := map[string]interface{}{
		"type":        "subscribe",
		"product_ids": c.ProductIDs,
		"channels":    c.Channels,
	}
	message, err := json.Marshal(subscription)
	if err != nil {
		log.Fatal("Error marshaling subscription:", err)
	}
	err = c.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Fatal("Error subscribing:", err)
	}
}

// Listen starts listening for messages from the WebSocket server
func (c *WebsocketClient) Listen() {
	defer c.Conn.Close()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		//log.Printf("\nReceived message: %s\n", message, "\n")
		c.MessageChannel <- string(message)
	}
}

// StartWebSocketClient starts the WebSocket client
func StartWebSocketClient(watchlist []string, messageChannel chan string) {
  fmt.Println("Starating Websocket Server for Coinbase")
	client := NewWebsocketClient(watchlist, "wss://ws-feed.pro.coinbase.com", nil, messageChannel)
	client.Connect("wss://ws-feed.pro.coinbase.com")
	client.Subscribe()
	go client.Listen()

	// Keep the main goroutine alive
	for {
		time.Sleep(1 * time.Second)
	}
}

