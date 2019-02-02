package main

import (
	"context"
	. "github.com/Fallensouls/Pandora/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	router := SetRouter()
	s := &http.Server{
		Addr:           ":" + Server.Port,
		Handler:        router,
		ReadTimeout:    Server.ReadTimeout * time.Second,
		WriteTimeout:   Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Panicf("Fail to start server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server......")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Panic("Server Shutdown:", err)
	}

	log.Printf("Server closed at %s", time.Now())
}
