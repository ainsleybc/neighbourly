package app_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

func TestSignUpUser(t *testing.T) {

	// connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
	})

	// close session on end test
	defer session.Close()

	// create the tables for test
	r.TableCreate("users").RunWrite(session)

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
	testRouter.RegisterHandler("user signup", SignUpUser)

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"user signup", ` +
		`"data":{"username":"david", "email":"david@david.com", "postcode":"wa12bj","password":"password"}}`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("users").Nth(0).Run(session)

	var row map[string]string
	var david string
	// res.One(&row) <- try and use this thing
	for res.Next(&row) {
		david = row["username"]
	}
	got2, want2 := david, "david"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	// read from socket
	var output Message
	conn.ReadJSON(&output)

	// write assertion
	got2, want2 = output.Name, "user created, logged in"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	r.TableDrop("users").Wait().Exec(session)
}
