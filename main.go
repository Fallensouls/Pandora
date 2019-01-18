package main

import (
	"Pandora/routers"
	"net/http"
	"time"
)

func main() {
	router, config := routers.SetRouter()
	s := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        router,
		ReadTimeout:    config.ReadTimeout  * time.Second,
		WriteTimeout:   config.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
