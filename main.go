package main

import (
	"app/config"
	pythonnodes "app/python_nodes"
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

	wg.Add(3)

	go func() {
		defer wg.Done()
		pythonnodes.RunPythonServer(config.GetPythonNodePort())
	}()

	go func() {
		defer wg.Done()
		server := &http.Server{
			Addr:           ":" + config.GetAppPort(),
			Handler:        router.AppRouter(),
			ReadTimeout:    300 * time.Second,
			WriteTimeout:   300 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		// Sử dụng ListenAndServeTLS để chạy server với HTTPS
		log.Fatalln(server.ListenAndServeTLS("keys/server.crt", "keys/server.key"))
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

		// Sử dụng ListenAndServeTLS cho socket server
		log.Fatalln(socker.ListenAndServeTLS("keys/server.crt", "keys/server.key"))
	}()

	go func() {
		defer wg.Done()
		rabbitmq.RunRabbitmq()
	}()

	wg.Wait()
}
