package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	externalAddr            = ":8080"
	externalReadTimeout     = time.Second * 2
	externalWriteTimeout    = time.Second * 2
	externalShutdownTimeout = time.Second * 10

	internalAddr            = ":8090"
	internalReadTimeout     = time.Second * 2
	internalWriteTimeout    = time.Second * 2
	internalShutdownTimeout = time.Second * 10
)

func main() {
	defer log.Println("service stopped")

	// External HTTP server

	externalMux := http.NewServeMux()

	externalMux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if _, err := rw.Write([]byte("Ahoy!")); err != nil {
			log.Printf("failed to write response: %+v", err)
		}
	})

	externalServer := &http.Server{
		Addr:         externalAddr,
		Handler:      externalMux,
		ReadTimeout:  externalReadTimeout,
		WriteTimeout: externalWriteTimeout,
	}

	go func() {
		if err := externalServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen and serve external HTTP server: %+v", err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), externalShutdownTimeout)
		defer cancel()

		if err := externalServer.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown external HTTP server: %+v", err)
		}
	}()

	// Internal HTTP server

	internalMux := http.NewServeMux()

	internalMux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		if _, err := rw.Write([]byte("healthy")); err != nil {
			log.Printf("failed to write health response: %+v", err)
		}
	})

	internalServer := &http.Server{
		Addr:         internalAddr,
		Handler:      internalMux,
		ReadTimeout:  internalReadTimeout,
		WriteTimeout: internalWriteTimeout,
	}

	go func() {
		if err := internalServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen and serve internal HTTP server: %+v", err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), internalShutdownTimeout)
		defer cancel()

		if err := internalServer.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown internal HTTP server: %+v", err)
		}
	}()

	// Waiting for shutdown

	log.Println("service is started")

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan

	log.Println("stopping service...")
}
