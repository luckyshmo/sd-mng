package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Messenger interface {
	Send(origin string, message []byte)
}

type WebSockets struct {
	sync.Mutex //! can cause really big performance issues on locked MU while write.
	connPool   map[string]*websocket.Conn
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
		ws.Lock()
		ws.connPool[origin] = c
		ws.Unlock()
		for {
			t, _, err := c.ReadMessage() //! can we concurrently read and write?
			if err != nil {
				fmt.Println(origin, ": read msg: ", err)
				ws.Lock()
				c.Close()
				delete(ws.connPool, origin)
				ws.Unlock()
				return
			}
			if t == websocket.CloseMessage {
				fmt.Println("close code for origin: ", origin)
				ws.Lock()
				c.Close()
				delete(ws.connPool, origin)
				ws.Unlock()
				return
			}
		}
	}
}

func (ws *WebSockets) Send(origin string, msg []byte) {
	ws.Lock()
	defer ws.Unlock()
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
