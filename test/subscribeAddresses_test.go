package app_test

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	"github.com/ainsleybc/neighbourly/db"
	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

func TestSubscribeAddresses(t *testing.T) {

	t.Parallel()

	db.CleanUp("subscribeAddresses")
	db.Setup("subscribeAddresses")
	defer db.CleanUp("subscribeAddresses")

	// connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "subscribeAddresses",
	})

	// close session on end test
	defer session.Close()

	// new router
	testRouter := NewRouter(session)

	// mock server thingy
	d := wstest.NewDialer(testRouter, nil)

	// // open websocket connection (skip error)
	conn, _, _ := d.Dial("ws://localhost:4000", nil)

	// register handler for addFeed message
	testRouter.RegisterHandler("feedAddress add", AddFeedAddress)
	testRouter.RegisterHandler("address subscribe", SubscribeAddress)

	// create test feed
	feed := map[string]string{"name": "General Assembly"}
	resp, _ := r.Table("feeds"). // create new feed
					Insert(feed).
					RunWrite(session)
	feedID := resp.GeneratedKeys[0]

	// user subscribing to their address feeds
	rawMessage := json.RawMessage(`{"name":"address subscribe", ` +
		`"data":{"feedId":"` + feedID + `"}}`)
	conn.WriteJSON(rawMessage)

	feedAddress := FeedAddress{
		Address: Address{
			StreetNumber: "1",
			StreetName:   "Makers Street",
			Postcode:     "AL74BD",
		},
		Feed: Feed{
			ID: feedID,
		}}

	r.Table("feedAddresses"). // create new feed
					Insert(feedAddress).
					RunWrite(session)

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// // readJSON from socket
	var output Message
	conn.ReadJSON(&output)

	// // write assertion
	got2, want2 := output.Name, "feedAddress add"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

}
