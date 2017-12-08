package app

import (
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
)

const (
	ChannelStop = iota
	UserStop
	MessageStop
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

func SubscribeFeed(client *Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(ChannelStop)
		cursor, _ := r.Table("feed").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		changeFeedHelper(cursor, "feed", client.send, stop)
	}()
}

func unsubscribeChannel(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string,
	send chan<- Message, stop <-chan bool) {
	change := make(chan r.ChangeResponse)
	cursor.Listen(change)
	for {
		eventName := ""
		var data interface{}
		select {
		case <-stop:
			cursor.Close()
			return
		case val := <-change:
			if val.NewValue != nil && val.OldValue == nil {
				eventName = changeEventName + " add"
				data = val.NewValue
				send <- Message{eventName, data}
			}
		}
	}
}
