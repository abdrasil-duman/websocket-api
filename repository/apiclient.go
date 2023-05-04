package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

type BestOrderBook struct {
	Ask Order `json:"ask"` //asks.Price > any bids.Price
	Bid Order `json:"bid"`
}

// Order struct
type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

type APIConnection struct {
	conn         *websocket.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewAPIConnection(conn *websocket.Conn) APIClient {
	return &APIConnection{
		conn:         conn,
		readTimeout:  10 * time.Second,
		writeTimeout: 10 * time.Second,
	}
}

// Connection connects to the AscendEX websocket API.
func (c *APIConnection) Connection() error {
	u := url.URL{
		Scheme: "wss",
		Host:   "ascendex.com",
		Path:   "/1/api/pro/v1/stream",
	}

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	c.conn = conn

	return nil
}

// Disconnect disconnects from the AscendEX websocket API.
func (c *APIConnection) Disconnect() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}
}

// SubscribeToChannel subscribes to the AscendEX BBO websocket channel for the given symbol.
func (c *APIConnection) SubscribeToChannel(symbol string) error {
	msg := map[string]interface{}{
		"method": "depth.subscribe",
		"params": []string{fmt.Sprintf("%s@depth20", symbol)},
		"id":     1,
	}

	err := c.conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("websocket write error: %w", err)
	}

	resp := make(map[string]interface{})
	err = c.conn.ReadJSON(&resp)
	if err != nil {
		return fmt.Errorf("websocket read error: %w", err)
	}

	if resp["id"].(float64) != 1 || resp["result"].(bool) != true {
		return errors.New("subscription failed")
	}

	return nil
}

// ReadMessagesFromChannel reads messages from the AscendEX websocket channel and writes them to the given channel.
func (c *APIConnection) ReadMessagesFromChannel(ch chan<- BestOrderBook) {
	for {
		err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			return
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			close(ch)
			return
		}

		var bbo BestOrderBook
		err = json.Unmarshal(msg, &bbo)
		if err != nil {
			close(ch)
			return
		}

		ch <- bbo
	}
}

// WriteMessagesToChannel writes messages to the AscendEX websocket channel.
func (c *APIConnection) WriteMessagesToChannel(msg []byte) error {
	err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return err
	}
	err = c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return fmt.Errorf("websocket write error: %w", err)
	}

	return nil
}
