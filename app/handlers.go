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
	feed := &Feed{
		Address: user.Postcode,
	}
	AddFeed(client, feed)
	cursor, _ := r.Table("feeds").
		Filter(r.Row.
			Field("address").
			Eq(user.Postcode)).
		Run(client.session)
	cursor.Next(&feed)
	user.DefaultFeed = feed.ID
	err := r.Table("users").
		Insert(user).
		Exec(client.session)
	if err != nil {
		client.send <- Message{Name: "signup unsuccesful"}
		return
	}
	client.user = user
	client.send <- Message{
		Name: "user created, logged in",
		Data: map[string]string{
			"email":       user.Email,
			"username":    user.Username,
			"defaultFeed": user.DefaultFeed,
		},
	}
}

func LoginUser(client *Client, data interface{}) {
	var login map[string]string
	var user User
	mapstructure.Decode(data, &login)
	cursor, _ := r.Table("users").
		Filter(r.Row.Field("email").
			Eq(login["email"])).
		Run(client.session)
	cursor.Next(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login["password"])); err != nil {
		client.send <- Message{Name: "incorrect credentials"}
		return
	}
	client.user = user
	client.send <- Message{
		Name: "login successful",
		Data: map[string]string{
			"email":       user.Email,
			"username":    user.Username,
			"defaultFeed": user.DefaultFeed,
		},
	}
}

func AddFeed(client *Client, data interface{}) {
	var feed Feed
	mapstructure.Decode(data, &feed)
	r.Table("feeds").
		Insert(feed).
		Exec(client.session)
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
		cursor, _ := r.Table("feeds").
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
