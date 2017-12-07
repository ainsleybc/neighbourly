package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

type Handler func(*Client, interface{})

func addFeed(client *Client, data interface{}) {
	var feed Feed
	mapstructure.Decode(data, &feed)
	go func() {
		r.Table("feed").
			Insert(feed).
			Exec(session)
	}()
}
