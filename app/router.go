package app

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Router struct {
	handlers map[string]Handler
	session  *r.Session
}
type Handler func(*Client, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *Router) RegisterHandler(msgName string, handler Handler) {
	r.handlers[msgName] = handler
}

func NewRouter(session *r.Session) *Router {
	return &Router{
		handlers: make(map[string]Handler),
		session:  session,
	}
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, _ := upgrader.Upgrade(w, r, nil)
	client := NewClient(socket, e.handlers, e.session)
	defer client.Close()
	// go client.Write()
	client.Read()

}
