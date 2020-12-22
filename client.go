package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client - struct
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send   chan []byte
	ID     int
	Status bool
	Remove bool
	Auth   bool
}

func (customer *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		customer.Conn.Close()
		customer.Status = false
	}()

	for {
		select {
		case message, ok := <-customer.Send:
			customer.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				customer.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := customer.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(customer.Send)
			for i := 0; i < n; i++ {
				writer.Write(newline)
				writer.Write(<-customer.Send)
			}

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			customer.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := customer.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (customer *Client) readPump() {

	defer func() {
		customer.Conn.Close()
		customer.Status = false
	}()

	customer.Conn.SetReadLimit(maxMessageSize)
	customer.Conn.SetReadDeadline(time.Now().Add(pongWait))
	customer.Conn.SetPongHandler(func(string) error { customer.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := customer.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		mes, err := json.Marshal(Letter{customer.ID, string(message)})
		if err != nil {
			log.Println(err)
		}

		MessageFromUsers <- mes
	}
}

// start - func for start methods of client
func (customer *Client) start() {
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go customer.writePump()
	go customer.readPump()
}

// GetID - return Client ID
// type int
func (customer Client) GetID() int {
	return customer.ID
}

// InitClient - (customer *websocket.Conn) *Client
func InitClient(customer *websocket.Conn) *Client {
	client := &Client{
		Conn: customer,
		Status: true,
		// Buffered channel of outbound messages.
		Send: make(chan []byte, maxMessageSize),
	}

	return client
}
