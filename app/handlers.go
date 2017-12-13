package app

import (
	"fmt"
	"strings"
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
	var address Address
	mapstructure.Decode(data, &user)
	mapstructure.Decode(data, &address)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	// default feed
	// TODO - create a separate method - insert if not exists, otherwise fetch - will currently error
	r.Table("addresses"). //create address
				Insert(address).
				RunWrite(client.session)

	feed := Feed{
		Name: strings.Join([]string{
			address.StreetNumber,
			address.StreetName,
			address.Postcode},
			" "),
	}

	resp, _ := r.Table("feeds"). // create new feed
					Insert(feed).
					RunWrite(client.session)

	res, _ := r.Table("feeds").
		Get(resp.GeneratedKeys[0]).
		Run(client.session)
	res.One(&feed)

	feedAddress := &FeedAddress{ // link the feed & address
		Feed:    feed,
		Address: address,
	}
	r.Table("feedAddresses").
		Insert(feedAddress).
		RunWrite(client.session)

	// assign default feed
	user.DefaultFeed = feed
	user.Address = address
	// insert new user
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
			"defaultFeed": user.DefaultFeed.ID,
		},
	}
}

func LoginUser(client *Client, data interface{}) {
	fmt.Println("User logged in")
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
	r.Table("feeds").
		Insert(feed).
		Exec(client.session)
}

func AddPost(client *Client, data interface{}) {
	fmt.Printf("\n\n\nData in AddPost function: %v\n\n\n", data)
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

func unsubscribeFeed(client *Client, data interface{}) {
	client.StopForKey(ChannelStop)
}

func SubscribePosts(client *Client, data interface{}) {

	go func() {
		fmt.Printf("\n\n\n\nSubscribe posts called in go \n\n\n")
		eventData := data.(map[string]interface{})
		val, _ := eventData["feedId"]
		fmt.Printf("%v\n\n\n", val)
		feedID, _ := val.(string)
		fmt.Printf("\n\nfeedID: %v\n\n", feedID)
		stop := client.NewStopChannel(MessageStop)
		cursor, _ := r.Table("posts").
			// OrderBy(r.OrderByOpts{Index: r.Desc("createdAt")}).
			Filter(r.Row.Field("feed").Eq(feedID)).
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		// fmt.Printf("\n\nCursor:%v\n\n\n", cursor)
		changeFeedHelper(cursor, "post", client.send, stop)
	}()
}

func unsubscribePosts(client *Client, data interface{}) {
	client.StopForKey(MessageStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string,
	send chan<- Message, stop <-chan bool) {
	fmt.Printf("\nIwas Called!!!!!\n")
	change := make(chan r.ChangeResponse)
	cursor.Listen(change)
	for {
		eventName := ""
		var data interface{}
		select {
		case <-stop:
			fmt.Printf("\n\nThe other case block!!v\n\n\n")

			cursor.Close()
			return
		case val := <-change:
			if val.NewValue != nil && val.OldValue == nil {
				fmt.Printf("Val: %v\n\n\n", val)
				eventName = changeEventName + " add"
				data = val.NewValue
				send <- Message{eventName, data}
			}
		}
	}
}
