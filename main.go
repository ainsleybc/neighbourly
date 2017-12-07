package main

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":4000", nil)
}

type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Router struct {
	// rules   map[string]Handler
	session *r.Session
}

func NewRouter(session *r.Session) *Router {
	return &Router{
		// rules:   make(map[string]Handler),
		session: session,
	}
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.Upgrade(w, r, nil)
	// client := NewClient(socket, e.FindHandler, e.session)
	// defer client.Close()
	// go client.Write()
	// client.Read()

}
