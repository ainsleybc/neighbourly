package app

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Handle func(*Client, Message)

type Client struct {
	socket       *websocket.Conn
	Handle       Handle
	session      *r.Session
	send         chan Message
	stopChannels map[int]chan bool
}

func NewClient(socket *websocket.Conn, handle Handle, session *r.Session) *Client {
	return &Client{
		socket:       socket,
		Handle:       handle,
		session:      session,
		send:         make(chan Message),
		stopChannels: make(map[int]chan bool),
	}
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
}

func (c *Client) Read() {
	var message Message
	for {
		if err := c.socket.ReadJSON(&message); err != nil {
			fmt.Printf("%v\n", err)
			break
		}
		c.Handle(c, message)
	}
	c.socket.Close()
}
