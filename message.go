package main

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type User struct {
	Id   string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}
