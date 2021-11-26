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
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// mux NewRouter
	r := mux.NewRouter()

	// product handlers
	ph := handlers.NewProducts(logger)

	// subrouter for GET methods
	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/product", ph.GetProducts)

	// subrouter for post method
	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/product", ph.CreateProduct)
	postRouter.Use(ph.MiddlewareProductValidator)

	// subrouter for put mehtod
	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/product/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidator)

	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           r,
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
