package play

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

func (c *Connection) SendMessage(messageType MessageType, data interface{}) {
	msg := Message{Type: messageType, Data: data}
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("error marshalling message: %v", err)
	}

	c.writeLock.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, bytes)
	c.writeLock.Unlock()
	if err != nil {
		log.Println("write:", err)
		return
	}
}
