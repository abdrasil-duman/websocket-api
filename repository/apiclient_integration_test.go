package repository

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestAPIClient(t *testing.T) {
	// Подключаемся к бирже с тестовым API ключом
	conn, _, err := websocket.DefaultDialer.Dial("wss://ascendex.com/123456", nil)
	assert.NoError(t, err)

	client := &APIConnection{
		conn: conn,
	}

	// Тестируем подключение к бирже
	err = client.Connection()
	assert.NoError(t, err)

	// Тестируем отключение от биржи
	client.Disconnect()
	_, _, err = conn.ReadMessage()
	assert.Error(t, err)

	// Тестируем подписку на канал BBO
	err = client.SubscribeToChannel("BTC_USDT")
	assert.NoError(t, err)

	// Тестируем чтение сообщений из канала
	ch := make(chan BestOrderBook)
	go client.ReadMessagesFromChannel(ch)
	book := <-ch
	assert.NotEmpty(t, book.Ask)
	assert.NotEmpty(t, book.Bid)
}
