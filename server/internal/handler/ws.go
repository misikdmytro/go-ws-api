package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/misikdmitriy/go-ws-api/internal/array"
	"github.com/misikdmitriy/go-ws-api/internal/client"
	wshandler "github.com/misikdmitriy/go-ws-api/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients client.WebSocketClientsPool
	m       sync.Mutex
)

func WebSocketConnect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go startClient(c, ws)
}

func startClient(ctx context.Context, ws *websocket.Conn) {
	cl := client.NewWebSocketClient(ws)

	m.Lock()
	clients = append(clients, cl)
	m.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("error: %v", err)
		}

		m.Lock()
		defer m.Unlock()
		clients = array.Except(clients, func(item client.WebSocketClient) bool { return item.Id() == cl.Id() })
		cl.Close()
	}()

	cl.Launch(ctx)
	wshandler.MemberJoin(clients, cl)

	for {
		select {
		case msg, ok := <-cl.Listen():
			if !ok {
				return
			} else {
				switch msg.Type {
				case "MESSAGE":
					wshandler.NewMessage(clients, cl, msg.Content["message"])
				}
			}
		case err := <-cl.Error():
			log.Printf("web socket error: %v", err)
		case <-cl.Done():
			wshandler.MemberLeave(clients, cl)
			return
		}
	}
}
