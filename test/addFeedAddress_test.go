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

func TestAddFeedAddress(t *testing.T) {

	t.Parallel()

	db.CleanUp("addFeedAddress")
	db.Setup("addFeedAddress")
	// defer db.CleanUp("addFeedAddress")

	// connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "addFeedAddress",
	})

	// close session on end test
	defer session.Close()

	// new router
	testRouter := NewRouter(session)

	// mock server thingy
	d := wstest.NewDialer(testRouter, nil)

	// open websocket connection (skip error)
	conn, _, _ := d.Dial("ws://localhost:4000", nil)

	// register handler for addFeedAddress message
	testRouter.RegisterHandler("feedAddress add", AddFeedAddress)

	// create test feed
	feed := map[string]string{"name": "General Assembly"}
	resp, _ := r.Table("feeds"). // create new feed
					Insert(feed).
					RunWrite(session)
	feedID := resp.GeneratedKeys[0]

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"feedAddress add", ` +
		`"data":{
			"address":{
				"streetNumber":"1",
				"streetName":"Makers Street",
				"postcode":"AL74BD"
			},
			"feed":{
				"ID":"` + feedID + `"
			}
			}}`)

	err := conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}
	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("feedAddresses").Nth(0).Run(session)

	var row map[string]interface{}
	// res.One(&row) <- try and use this thing
	res.Next(&row)
	got1, want1 := row["feed"].(string), feedID
	if got1 != want1 {
		t.Errorf("got: %v, want: %v", got1, want1)
	}

}
