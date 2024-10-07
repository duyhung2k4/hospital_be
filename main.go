package main

import (
	"app/config"
	"app/router"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Cấu hình HTTPS
		httpsServer := http.Server{
			Addr:           ":" + config.GetAppPort(),
			Handler:        router.AppRouter(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		// Đường dẫn đến chứng chỉ và khóa riêng
		certFile := "keys/server.crt"
		keyFile := "keys/server.key"

		// Listen và Serve trên HTTPS
		log.Println("HTTPS server is running on port 10000")
		log.Fatalln(httpsServer.ListenAndServeTLS(certFile, keyFile))
	}()

	wg.Wait()
}
