package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func (server Server) Start(ctxCancel context.CancelFunc) {

	log.Println(fmt.Sprintf("start server listen at [%s]", server.srv.Addr))

	if err := server.srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("server closed under request")
		} else {
			log.Println("server closed unexpect:", err)
		}
	}

	ctxCancel()

}

func (server Server) Close() error {
	return server.srv.Close()
}

func NewServer(config Config) *Server {

	log.Println("initializing server")

	srv := &http.Server{
		Addr:    config.Addr(),
		Handler: router(config.PublicPath()),
	}

	server := &Server{
		srv: srv,
	}

	return server

}
