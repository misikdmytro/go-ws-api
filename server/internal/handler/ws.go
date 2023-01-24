package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/misikdmitriy/go-ws-api/internal/array"
	"github.com/misikdmitriy/go-ws-api/internal/client"
	"github.com/misikdmitriy/go-ws-api/internal/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients []client.WebSocketClient = make([]client.WebSocketClient, 0)

func WebSocketConnect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go func() {
		cl := client.NewWebSocketClient(ws)
		clients = append(clients, cl)
		defer func() {
			clients = array.Except(clients, func(item client.WebSocketClient) bool { return item.Id() == cl.Id() })
			cl.Close()
		}()

		cl.Launch(c)

		for {
			select {
			case msg, ok := <-cl.Listen():
				if !ok {
					return
				} else {
					if msg.Type == "NEW_CLIENT" {
						cl.Write(model.WebSocketMessage{
							Type: "ID_ASSIGNED",
							Content: map[string]string{
								"id": cl.Id(),
							},
						})

						array.ForEach(
							array.Except(clients, func(item client.WebSocketClient) bool { return item.Id() == cl.Id() }),
							func(item client.WebSocketClient) {
								item.Write(model.WebSocketMessage{
									Type: "MEMBER_JOIN",
									Content: map[string]string{
										"id": cl.Id(),
									},
								})
							},
						)
					}
				}
			case err := <-cl.Error():
				log.Printf("web socket error: %v", err)
			case <-cl.Done():
				array.ForEach(
					array.Except(clients, func(item client.WebSocketClient) bool { return item.Id() == cl.Id() }),
					func(item client.WebSocketClient) {
						item.Write(model.WebSocketMessage{
							Type: "MEMBER_LEAVE",
							Content: map[string]string{
								"id": cl.Id(),
							},
						})
					},
				)
				return
			}
		}
	}()
}
