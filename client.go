package main

import "github.com/gorilla/websocket"

// r "github.com/dancannon/gorethink"
// "github.com/gorilla/websocket"
// "log"

type Client struct {
	send   chan Message
	socket *websocket.Conn
	// findHandler FindHandler
	// session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}
