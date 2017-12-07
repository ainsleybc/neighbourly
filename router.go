package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Router struct {
	rules map[string]Handler
}
type Handle func(*Client, Message)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *Router) RegisterHandler(msgName string, handler Handler) {
	r.rules[msgName] = handler
}

func (r *Router) Handle(client *Client, msg Message) {
	handler := r.rules[msg.Name]
	handler(client, msg.Data)
}

func NewRouter() *Router {
	return &Router{
		rules: make(map[string]Handler),
	}
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, _ := upgrader.Upgrade(w, r, nil)
	client := NewClient(socket, e.Handle)
	defer client.Close()
	// go client.Write()
	client.Read()

}
