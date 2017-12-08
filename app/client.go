package app

import (
	"fmt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
)

type Handle func(*Client, Message)

type Client struct {
	socket       *websocket.Conn
	handlers     map[string]Handler
	session      *r.Session
	send         chan Message
	stopChannels map[int]chan bool
}

func NewClient(params ...interface{}) *Client {
	return &Client{
		socket:       params[0].(*websocket.Conn),
		handlers:     params[1].(map[string]Handler),
		session:      params[2].(*r.Session),
		send:         make(chan Message),
		stopChannels: make(map[int]chan bool),
	}
}

func (c *Client) Handle(msg Message) {
	handler := c.handlers[msg.Name]
	handler(c, msg.Data)
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
		c.Handle(message)
	}
	c.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			fmt.Printf("%v", err)
			break
		}
	}
	client.socket.Close()
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}
