package client

import (
	"context"
	"errors"
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
	msgs chan model.WebSocketMessage
	err  chan error
	done chan interface{}
	m    sync.Mutex
	o    sync.Once
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

type WebSocketClientsPool []WebSocketClient

func NewWebSocketClient(ws *websocket.Conn) WebSocketClient {
	return &webSocketClient{
		id:   uuid.NewString(),
		ws:   ws,
		msgs: make(chan model.WebSocketMessage),
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

	c.o.Do(func() { go c.launch(ctx) })
}

func (c *webSocketClient) launch(ctx context.Context) {
	var wg sync.WaitGroup

	cancellation, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.write(websocket.CloseMessage)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.read(cancellation)
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.ping(cancellation)
		cancel()
	}()

	wg.Wait()
	c.done <- struct{}{}
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
				c.handleError(err)
				return
			}

			c.msgs <- msg
		}
	}
}

func (c *webSocketClient) ping(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.write(websocket.PingMessage)
		case <-ctx.Done():
			return
		}
	}
}

func (c *webSocketClient) write(messageType int) {
	c.m.Lock()
	defer c.m.Unlock()

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := c.ws.WriteMessage(messageType, nil); err != nil {
		c.handleError(err)
	}
}

func (c *webSocketClient) handleError(err error) {
	if _, ok := err.(*websocket.CloseError); ok {
		return
	}

	if errors.Is(err, websocket.ErrCloseSent) {
		return
	}

	c.err <- err
}

func (c *webSocketClient) Listen() <-chan model.WebSocketMessage {
	return c.msgs
}

func (c *webSocketClient) Done() <-chan interface{} {
	return c.done
}

func (c *webSocketClient) Error() <-chan error {
	return c.err
}

func (c *webSocketClient) Write(m model.WebSocketMessage) error {
	c.m.Lock()
	defer c.m.Unlock()

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(m)
}

func (c *webSocketClient) Close() {
	c.ws.Close()
	close(c.msgs)
}
