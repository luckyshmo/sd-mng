package main

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
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

func (ws *WebSockets) ProgressHandler(conn *websocket.Conn) {
	ws.connPool[conn.Config().Origin.String()] = conn
	for {
		time.Sleep(time.Millisecond * 10000)
		err := websocket.Message.Send(conn, Message{"1111", "kek.file", 75})
		if err != nil {
			fmt.Println("ERROR sending fake message: ", err)
		}
	}
}

func (ws *WebSockets) Send(origin string, msg []byte) {
	if conn, ok := ws.connPool[origin]; ok {
		err := websocket.Message.Send(conn, string(msg))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		return
	}
	fmt.Println(origin, " not connected")
}
