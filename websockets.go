package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Messenger interface {
	Send(origin string, message []byte)
}

type WebSockets struct {
	connPool map[string]*websocket.Conn
}

func NewWebSockets() *WebSockets {
	return &WebSockets{
		connPool: map[string]*websocket.Conn{},
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

func (ws *WebSockets) ProgressHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		origin := r.Header.Get("Origin")
		ws.connPool[origin] = c
		for {
			t, _, err := c.ReadMessage()
			if err != nil {
				fmt.Println(origin, ": read msg: ", err)
				c.Close()
				delete(ws.connPool, origin)
				return
			}
			if t == websocket.CloseMessage {
				fmt.Println("close code for origin: ", origin)
				c.Close()
				delete(ws.connPool, origin)
				return
			}
		}
	}
}

func (ws *WebSockets) Send(origin string, msg []byte) {
	if conn, ok := ws.connPool[origin]; ok {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			delete(ws.connPool, origin)
		}
		return
	}
	fmt.Println(origin, " not connected")
}
