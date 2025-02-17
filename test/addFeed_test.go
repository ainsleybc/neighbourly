package app_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	"github.com/ainsleybc/neighbourly/db"
	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

func TestAddFeed(t *testing.T) {

	t.Parallel()

	db.CleanUp("addFeed")
	db.Setup("addFeed")
	defer db.CleanUp("addFeed")

	// connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "addFeed",
	})

	// close session on end test
	defer session.Close()

	// new router
	testRouter := NewRouter(session)

	// mock server thingy
	d := wstest.NewDialer(testRouter, nil)

	// open websocket connection (skip error)
	conn, resp, err := d.Dial("ws://localhost:4000", nil)

	// assertion 1 (check websocket upgrade connection status)
	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %v, want: %v", got, want)
	}

	// register handler for addFeed message
	testRouter.RegisterHandler("feed add", AddFeed)

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"feed add", ` +
		`"data":{"Name":"Makers Academy"}}`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}
	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("feeds").Nth(0).Run(session)

	var row map[string]string
	var david string
	// res.One(&row) <- try and use this thing
	for res.Next(&row) {
		david = row["name"]
	}
	got2, want2 := david, "Makers Academy"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

}
