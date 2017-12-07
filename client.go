package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	socket       *websocket.Conn
	Handle       Handle
	send         chan Message
	stopChannels map[int]chan bool
}

func NewClient(socket *websocket.Conn, handle Handle) *Client {
	return &Client{
		socket:       socket,
		Handle:       handle,
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

func (client *Client) Read() {
	var message Message
	var data []byte
	for {
		if err := client.socket.ReadJSON(&data); err != nil {
			break
		}
		json.Unmarshal(data, &message)
		client.Handle(client, message)
	}
	client.socket.Close()
}
