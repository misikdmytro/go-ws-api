package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/go-ws-api/internal/handler"
)

func NewHttpServer() *http.Server {
	r := gin.Default()

	r.GET("/ws", handler.WebSocketConnect)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
