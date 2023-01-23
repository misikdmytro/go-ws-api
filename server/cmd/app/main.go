package main

import "github.com/misikdmitriy/go-ws-api/internal/server"

func main() {
	s := server.NewHttpServer()
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
