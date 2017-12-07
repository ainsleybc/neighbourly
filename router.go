package main

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Router struct {
	// rules   map[string]Handler
	session *r.Session
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *Router) Handle(msgName string, handler Handler) {
	// r.rules[msgName] = handler
}

// func (r *Router) FindHandler(msgName string) (Handler, bool) {
// 	// handler, found := r.rules[msgName]
// 	return handler, found
// }

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
