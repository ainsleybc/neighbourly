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

func TestLoginUser(t *testing.T) {

	t.Parallel()

	db.Setup("loginUser")
	defer db.CleanUp("loginUser")

	// connect to rethinkDB
	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "loginUser",
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

	// register handlers
	testRouter.RegisterHandler("user signup", SignUpUser)
	testRouter.RegisterHandler("user login", LoginUser)

	// sign up a user
	rawMessage := json.RawMessage(`{"name":"user signup", ` +
		`"data":{"username":"david", "email":"david@david.com", "postcode":"wa12bj","password":"password"}}`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}

	conn.Close() // end session

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	d = wstest.NewDialer(testRouter, nil)
	conn, resp, err = d.Dial("ws://localhost:4000", nil) // start new ws session

	// login a user
	rawMessage = json.RawMessage(`{"name":"user login", ` +
		`"data":{"email":"david@david.com","password":"password"}}`)

	// // stuff message
	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}

	// read from socket
	var output Message
	conn.ReadJSON(&output)
	// write assertion

	got2, want2 := output.Name, "login successful"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

}
