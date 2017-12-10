package app

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

func TestSubscribePost(t *testing.T) {

	t.Parallel()

	db.Setup("subscribePost")
	defer db.CleanUp("subscribePost")

	// // connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "subscribePost",
	})

	// close session on end test
	defer session.Close()

	// new router
	testRouter := NewRouter(session)

	// mock server thingy
	d := wstest.NewDialer(testRouter, nil)

	// open websocket connection (skip error)
	conn, resp, _ := d.Dial("ws://localhost:4000", nil)

	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %q, want: %q", got, want)
	}

	// register handler for post subscribe message
	testRouter.RegisterHandler("post subscribe", SubscribePosts)

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"post subscribe","data":{"feedId": "123hhsj111"}}`)
	conn.WriteJSON(rawMessage)

	// create post & add to database
	post := &Post{
		Name:      "Jon",
		CreatedAt: time.Now(),
		Text:      "Subscribing!",
		FeedID:    "123hhsj111",
	}
	r.Table("posts").Insert(post).RunWrite(session)

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// readJSON from socket
	var output Message
	conn.ReadJSON(&output)

	// write assertion
	got2, want2 := output.Name, "post add"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

}
