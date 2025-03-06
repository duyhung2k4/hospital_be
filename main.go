package main

import (
	"app/config"
	"app/rabbitmq"
	"app/router"
	"app/socket"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		server := &http.Server{
			Addr:           ":" + config.GetAppPort(),
			Handler:        router.AppRouter(),
			ReadTimeout:    300 * time.Second,
			WriteTimeout:   300 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatalln(server.ListenAndServe())
	}()

	go func() {
		defer wg.Done()
		socker := &http.Server{
			Addr:           ":" + config.GetSocketPort(),
			Handler:        socket.ServerSocker(),
			ReadTimeout:    300 * time.Second,
			WriteTimeout:   300 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		log.Fatalln(socker.ListenAndServe())
	}()

	go func() {
		defer wg.Done()
		rabbitmq.RunRabbitmq()
	}()

	wg.Wait()
}
