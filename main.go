package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajaymahar/microservices/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(logger)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           sm,
		TLSConfig:         &tls.Config{},
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       3 * time.Second,
		ErrorLog:          logger,
	}

	// NOTE: Server to start listening requests
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			s.ErrorLog.Fatal(err.Error())
		}
	}()

	// NOTE: chan to keep listening for os sign
	sigChan := make(chan os.Signal)

	// NOTE: any os action interrupt or kill will will notify the sigChan
	// signal.Notify(sigChan, os.Interrupt)
	// signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGTERM)

	// NOTE: waiting to any os signal interruption to do gracefull shoutdown
	sig := <-sigChan
	logger.Println("Recieved signal", sig)

	// NOTE: context with timeout of 30 sec, to do gracefull shoutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO: think about the shoutdown
	err := s.Shutdown(ctx)
	if err != nil {
		s.ErrorLog.Println("Error while shoutting down the server", err.Error())
	}
}
