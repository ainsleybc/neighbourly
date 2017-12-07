package main

import "github.com/gorilla/websocket"

type Client struct {
	send   chan Message
	socket *websocket.Conn
	// findHandler FindHandler
	// session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}
