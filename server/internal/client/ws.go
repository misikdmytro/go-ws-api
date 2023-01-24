package client

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/misikdmitriy/go-ws-api/internal/model"
)

const (
	maxMessageSize = 1024
	writeWait      = 5 * time.Second
	pingPeriod     = 5 * time.Second
	pongWait       = 10 * time.Second
)

type webSocketClient struct {
	id   string
	ws   *websocket.Conn
	buff chan model.WebSocketMessage
	i    chan model.WebSocketMessage
	err  chan error
	done chan interface{}
}

type WebSocketClient interface {
	Id() string
	Launch(ctx context.Context)
	Write(m model.WebSocketMessage) error
	Close()
	Listen() <-chan model.WebSocketMessage
	Done() <-chan interface{}
	Error() <-chan error
}

func NewWebSocketClient(ws *websocket.Conn) WebSocketClient {
	return &webSocketClient{
		id:   uuid.NewString(),
		ws:   ws,
		buff: make(chan model.WebSocketMessage),
		i:    make(chan model.WebSocketMessage),
		err:  make(chan error),
		done: make(chan interface{}),
	}
}

func (c *webSocketClient) Id() string {
	return c.id
}

func (c *webSocketClient) Launch(ctx context.Context) {
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.read(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.write(ctx)
		}()

		wg.Wait()
		c.done <- struct{}{}
	}()
}

func (c *webSocketClient) read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg model.WebSocketMessage
			err := c.ws.ReadJSON(&msg)
			if err != nil {
				c.err <- err
				return
			}

			c.buff <- msg
		}
	}
}

func (c *webSocketClient) write(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.buff:
			if !ok {
				if err := c.ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
					c.err <- err
					return
				}
			} else {
				c.i <- msg
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.err <- err
				return
			}
		case <-ctx.Done():
			if err := c.ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
				c.err <- err
			}

			return
		}
	}
}

func (c *webSocketClient) Listen() <-chan model.WebSocketMessage {
	return c.i
}

func (c *webSocketClient) Done() <-chan interface{} {
	return c.done
}

func (c *webSocketClient) Error() <-chan error {
	return c.err
}
func (c *webSocketClient) Write(m model.WebSocketMessage) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(m)
}

func (c *webSocketClient) Close() {
	c.ws.Close()

	close(c.buff)
	close(c.i)
}
