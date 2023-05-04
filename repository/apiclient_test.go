package repository

import (
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
)

type mockAPIClient struct {
	conn     *websocket.Conn
	symbol   string
	isCalled bool
}

func (c *mockAPIClient) Connection() error {
	if c.conn == nil {
		return fmt.Errorf("connection failed")
	}
	return nil
}

func (c *mockAPIClient) Disconnect() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return
		}
	}
}

func (c *mockAPIClient) SubscribeToChannel(symbol string) error {
	if symbol == "" {
		return fmt.Errorf("symbol is empty")
	}
	c.symbol = symbol
	return nil
}

func (c *mockAPIClient) ReadMessagesFromChannel(ch chan<- BestOrderBook) {
	if c.conn == nil {
		return
	}
	for {
		if c.isCalled {
			break
		}
	}
	ch <- BestOrderBook{
		Ask: Order{
			Amount: 1.0,
			Price:  100.0,
		},
		Bid: Order{
			Amount: 1.0,
			Price:  99.0,
		},
	}
}

func (c *mockAPIClient) WriteMessagesToChannel() {
	if c.conn == nil {
		return
	}
	c.isCalled = true
}

func TestMockAPIClient(t *testing.T) {
	c := &mockAPIClient{}
	err := c.Connection()
	if err == nil {
		t.Errorf("Connection() should return an error when conn is nil")
	}
	c.conn = &websocket.Conn{}
	err = c.Connection()
	if err != nil {
		t.Errorf("Connection() should not return an error when conn is not nil")
	}
	c.Disconnect()
	if c.conn != nil {
		t.Errorf("Disconnect() should close the connection")
	}
	err = c.SubscribeToChannel("")
	if err == nil {
		t.Errorf("SubscribeToChannel() should return an error when symbol is empty")
	}
	err = c.SubscribeToChannel("USDT_BTC")
	if err != nil {
		t.Errorf("SubscribeToChannel() should not return an error when symbol is not empty")
	}
	ch := make(chan BestOrderBook)
	go c.ReadMessagesFromChannel(ch)
	bbo := <-ch
	if bbo.Ask.Price != 100.0 || bbo.Bid.Price != 99.0 {
		t.Errorf("ReadMessagesFromChannel() should return valid data")
	}
	c.WriteMessagesToChannel()
	if !c.isCalled {
		t.Errorf("WriteMessagesToChannel() should set isCalled to true")
	}
}
