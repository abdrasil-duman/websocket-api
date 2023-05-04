<h1>AscendEX Websocket API</h1>
This is a Golang project for connecting to the AscendEX websocket API. It allows users to subscribe to the BBO (Best Bid Offer) channel for a given symbol and read messages from it.

Installation
To install the package, run:

```bash
go get github.com/abdrasil-duman/websocket-api
```

### Usage

To use the package, first create a new APIConnection by calling NewAPIConnection with a websocket connection. Then,
connect to the websocket API by calling Connection. You can then subscribe to the BBO channel for a given symbol by
calling SubscribeToChannel, and start reading messages from the channel by calling ReadMessagesFromChannel.

```go
package main

import (
"fmt"
"github.com/gorilla/websocket"
)

func main() {
// Create websocket connection
conn, _, err := websocket.DefaultDialer.Dial("wss://ascendex.com/1/api/pro/v1/stream", nil)
if err != nil {
fmt.Println("websocket dial error:", err)
return
}
// Create APIConnection
apiConn := repository.NewAPIConnection(conn)

// Connect to websocket API
err = apiConn.Connection()
if err != nil {
  fmt.Println("websocket connection error:", err)
  return
	}

// Subscribe to BBO channel for BTC/USDT symbol
  err = apiConn.SubscribeToChannel("BTCUSDT")
  if err != nil {
	fmt.Println("subscription error:", err)
 	return
  	}

// Read messages from channel
ch := make(chan repository.BestOrderBook)
go apiConn.ReadMessagesFromChannel(ch)

for bbo := range ch {
		// Do something with bbo
	}

}
```

### APIConnection Methods

- NewAPIConnection(conn *websocket.Conn) APIClient: Creates a new APIConnection with the given websocket connection.
- Connection() error: Connects to the AscendEX websocket API.
- Disconnect(): Disconnects from the AscendEX websocket API.
- SubscribeToChannel(symbol string) error: Subscribes to the BBO channel for the given symbol.
- ReadMessagesFromChannel(ch chan<- BestOrderBook): Reads messages from the BBO channel and writes them to the given
  channel.
- WriteMessagesToChannel(msg []byte) error: Writes messages to the AscendEX websocket channel.