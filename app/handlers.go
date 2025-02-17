package app

import (
	"strings"
	"time"

	"github.com/ainsleybc/neighbourly/db"
	r "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

const (
	ChannelStop = iota
	UserStop
	MessageStop
	FeedAddressStop
)

func SignUpUser(client *Client, data interface{}) {
	var user User
	var address Address
	mapstructure.Decode(data, &user)
	mapstructure.Decode(data, &address)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	addressPk := []string{
		address.StreetNumber,
		address.StreetName,
		address.Postcode,
	}

	feedName := strings.Join(addressPk, " ")

	feed := Feed{
		Name:           feedName,
		AddressDefault: true,
	}

	cursor, _ := db.GetDefaultFeedByAddress(client.session, addressPk)
	var row map[string]interface{}
	cursor.Next(&row)
	defaultFeed := row["feed"]

	if defaultFeed != nil {
		feed.ID = defaultFeed.(string)
	} else {
		// create new feed & address
		db.InsertAddress(client.session, address)
		resp, _ := db.InsertFeed(client.session, feed)
		feed.ID = resp.GeneratedKeys[0]
	}

	feedAddress := &FeedAddress{ // link the feed & address
		Feed:    feed,
		Address: address,
	}
	db.InsertFeedAddress(client.session, feedAddress)

	// assign default feed
	user.DefaultFeed = feed
	user.Address = address

	// insert new user
	_, err := db.InsertUser(client.session, user)

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
			"defaultFeed": user.DefaultFeed.ID,
		},
	}
}

func LoginUser(client *Client, data interface{}) {
	var login map[string]string
	var user User
	mapstructure.Decode(data, &login)
	cursor, _ := r.Table("users").
		Get(login["email"]).
		Merge(func(p r.Term) interface{} {
			return map[string]interface{}{
				"defaultFeed": r.Table("feeds").Get(p.Field("defaultFeed")),
			}
		}).
		Merge(func(p r.Term) interface{} {
			return map[string]interface{}{
				"address": r.Table("addresses").Get(p.Field("address")),
			}
		}).
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
			"defaultFeed": user.DefaultFeed.ID,
		},
	}
}

func AddFeed(client *Client, data interface{}) {
	var feed Feed
	mapstructure.Decode(data, &feed)
	resp, _ := r.Table("feeds"). // create new feed
					Insert(feed).
					RunWrite(client.session)

	feed.ID = resp.GeneratedKeys[0]

	feedAddress := &FeedAddress{ // link the feed & address
		Feed:    feed,
		Address: client.user.Address,
	}
	AddFeedAddress(client, feedAddress)
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
	address := []string{
		client.user.Address.StreetNumber,
		client.user.Address.StreetName,
		client.user.Address.Postcode,
	}
	go func() {
		stop := client.NewStopChannel(ChannelStop)
		cursor, _ := r.Table("feedAddresses").
			GetAllByIndex("address", address).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Map(r.Table("feeds"), func(res r.Term, feed r.Term) interface{} {
				return res.Merge(func(row r.Term) map[string]interface{} {
					return map[string]interface{}{
						"new_val": r.Table("feeds").Get(row.Field("new_val").Field("feed")),
					}
				})
			}).
			Run(client.session)
		changeFeedHelper(cursor, "feed", client.send, stop)
	}()
}

func UnsubscribeFeed(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

func SubscribePosts(client *Client, data interface{}) {
	go func() {
		eventData := data.(map[string]interface{})
		val, _ := eventData["feedId"]
		feedID, _ := val.(string)
		stop := client.NewStopChannel(MessageStop)
		cursor, _ := r.Table("posts").
			Filter(r.Row.Field("feed").Eq(feedID)).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		changeFeedHelper(cursor, "post", client.send, stop)
	}()
}

func UnsubscribePosts(client *Client, data interface{}) {
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
			if val.NewValue != nil {
				eventName = changeEventName + " add"
				data = val.NewValue
				send <- Message{eventName, data}
			}
		}
	}
}

func AddFeedAddress(client *Client, data interface{}) {
	var feedAddress FeedAddress
	mapstructure.Decode(data, &feedAddress)
	go func() {
		r.Table("feedAddresses").
			Insert(feedAddress).
			Exec(client.session)
	}()
}

func SubscribeAddress(client *Client, data interface{}) {
	go func() {
		eventData := data.(map[string]interface{})
		val, _ := eventData["feedId"]
		feedID, _ := val.(string)
		stop := client.NewStopChannel(FeedAddressStop)
		cursor, _ := r.Table("feedAddresses").
			Filter(r.Row.Field("feed").Eq(feedID)).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		changeFeedHelper(cursor, "feedAddress", client.send, stop)
	}()
}

func UnsubscribeAddress(client *Client, data interface{}) {
	client.StopForKey(FeedAddressStop)
}
