package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	container := bootstrap()
	http.HandleFunc("/ping", container.httpApp.Ping)
	http.HandleFunc("/chat", container.httpApp.CreateRoom)
	http.HandleFunc("/chat/join", container.websocketApp.Handle)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("Cannot start server - %v", err)
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	<-terminate
	log.Printf("Shutting down server...")

	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	container.websocketApp.Shutdown(ctx)
}
