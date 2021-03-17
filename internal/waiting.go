package internal

import (
	"context"
	"log"
	"os"
	"os/signal"
)

type Server interface {
	Start(ctxCancel context.CancelFunc)
	Close() error
}

type Redis interface {
	Close() error
}

func Waiting(server Server, redis Redis) {

	ctx, ctxCancel := context.WithCancel(context.Background())

	go server.Start(ctxCancel)

	log.Println("initializing waiting")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Println("Server Close:", err)
		}
		if err := redis.Close(); err != nil {
			log.Println("Redis Close:", err)
		}
	case <-ctx.Done():
		log.Println("shutdown by context")
	}

}
