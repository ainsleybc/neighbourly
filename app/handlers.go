package app

import (
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

const (
	ChannelStop = iota
	UserStop
	MessageStop
)

func SignUpUser(client *Client, data interface{}) {
	var user User
	mapstructure.Decode(data, &user)
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	r.Table("users").
		Insert(user).
		Exec(client.session)
	client.user = user
}

func AddFeed(client *Client, data interface{}) {
	var feed Feed
	mapstructure.Decode(data, &feed)
	go func() {
		r.Table("feed").
			Insert(feed).
			Exec(client.session)
	}()
}

func AddPost(client *Client, data interface{}) {
	var post Post
	mapstructure.Decode(data, &post)
	go func() {
		post.CreatedAt = time.Now()
		r.Table("posts").
			Insert(post).
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

func unsubscribeFeed(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

func SubscribePosts(client *Client, data interface{}) {
	go func() {
		eventData := data.(map[string]interface{})
		val, _ := eventData["feedId"]
		feedID, _ := val.(string)
		stop := client.NewStopChannel(MessageStop)
		cursor, _ := r.Table("posts").
			// OrderBy(r.OrderByOpts{r.Desc("createdAt")}).
			Filter(r.Row.Field("feedId").Eq(feedID)).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		changeFeedHelper(cursor, "post", client.send, stop)
	}()
}

func unsubscribePosts(client *Client, data interface{}) {
	client.StopForKey(MessageStop)
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
