package client

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/misikdmitriy/go-ws-api/internal/model"
)

const (
	maxMessageSize = 1024
	writeWait      = 5 * time.Second
	pingPeriod     = 55 * time.Second
	pongWait       = 60 * time.Second
)

type webSocketClient struct {
	ws *websocket.Conn
	m  chan model.WebSocketMessage
}

type WebSocketClient interface {
}

func NewWebSocketClient(ws *websocket.Conn) WebSocketClient {
	c := &webSocketClient{
		ws: ws,
		m:  make(chan model.WebSocketMessage),
	}

	go c.Read()
	go c.Write()

	return c
}

func (c *webSocketClient) Read() {
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg model.WebSocketMessage
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error on message read: %v", err)
			return
		}

		c.m <- msg
	}
}

func (c *webSocketClient) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	for {
		select {
		case msg, ok := <-c.m:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
			} else {
				if err := c.ws.WriteJSON(msg); err != nil {
					log.Printf("error on message write: %v", err)
					return
				}
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error on message write: %v", err)
				return
			}
		}
	}
}

func (c *webSocketClient) Close() {
	c.ws.Close()
	close(c.m)
}
