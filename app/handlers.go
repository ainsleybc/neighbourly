package app

import (
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

func AddFeed(client *Client, data interface{}) {
	var feed Feed
	mapstructure.Decode(data, &feed)
	go func() {
		r.Table("feed").
			Insert(feed).
			Exec(client.session)
	}()
}
