package controller

import (
	"app/config"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type socketController struct {
	upgrader       *websocket.Upgrader
	mutexSocket    *sync.Mutex
	mapSocket      map[string]*websocket.Conn
	mapSocketEvent map[string]map[string]*websocket.Conn
}

type SocketController interface {
	AuthSocket(w http.ResponseWriter, r *http.Request)
	FaceLoginSocket(w http.ResponseWriter, r *http.Request)
}

func (c *socketController) AuthSocket(w http.ResponseWriter, r *http.Request) {
	// check auth with uuid
	query := r.URL.Query()
	uuid := query.Get("uuid")
	if uuid == "" {
		badRequest(w, r, errors.New("uuid not found"))
		return
	}

	// create connect
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	//connect -> map socket
	c.mutexSocket.Lock()
	c.mapSocket[uuid] = conn
	c.mutexSocket.Unlock()

	// listen connect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
	}

	log.Printf("Disconnect")
}

func (c *socketController) FaceLoginSocket(w http.ResponseWriter, r *http.Request) {
	// check auth with uuid
	query := r.URL.Query()
	uuid := query.Get("uuid")
	if uuid == "" {
		badRequest(w, r, errors.New("uuid not found"))
		return
	}

	// create connect
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	//connect -> map socket
	c.mutexSocket.Lock()
	c.mapSocket[uuid] = conn
	c.mutexSocket.Unlock()

	// listen connect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
	}

	log.Printf("Disconnect")
}

func NewSocketController() SocketController {
	return &socketController{
		mutexSocket:    new(sync.Mutex),
		upgrader:       config.GetUpgraderSocket(),
		mapSocket:      config.GetMapSocket(),
		mapSocketEvent: config.GetSocketEvent(),
	}
}
