package app

import (
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Router struct {
	rules   map[string]Handler
	session *r.Session
}
type Handle func(*Client, Message)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// func (r *Router) RegisterHandler(msgName string, handler Handler) {
// 	r.rules[msgName] = handler
// }

func (r *Router) Handle(client *Client, msg Message) {
	handler := r.rules[msg.Name]
	handler(client, msg.Data)
}

func NewRouter(session *r.Session) *Router {
	return &Router{
		rules: map[string]Handler{
			"feed add": addFeed,
		},
		session: session,
	}
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, _ := upgrader.Upgrade(w, r, nil)
	client := NewClient(socket, e.Handle, e.session)
	defer client.Close()
	// go client.Write()
	client.Read()

}
