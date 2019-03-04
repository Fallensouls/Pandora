package main

import (
	"context"
	. "github.com/Fallensouls/Pandora/conf"
	"github.com/Fallensouls/Pandora/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := routers.SetRouter()
	server := &http.Server{
		Addr:           ":" + Config.Port,
		Handler:        router,
		ReadTimeout:    Config.ReadTimeout * time.Second,
		WriteTimeout:   Config.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("Fail to start server: %s", err)
		}
		log.Println("--------------Welcome to Pandora---------------")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server......")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Panic("Server Shutdown:", err)
	}

	log.Printf("Server closed at %s", time.Now())
}
