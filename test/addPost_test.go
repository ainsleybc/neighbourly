package app

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

func TestMain(m *testing.M) {
	// test set up
	result := m.Run()
	// test tear down
	os.Exit(result)
}

func TestAddPost(t *testing.T) {

	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
	})

	// close session on end test
	defer session.Close()

	// create the tables for test
	r.TableCreate("posts").RunWrite(session)

	// new router
	testPostRouter := NewRouter(session)

	// mock server thingy
	postDialer := wstest.NewDialer(testPostRouter, nil)

	// open websocket connection (skip error)
	conn, resp, err := postDialer.Dial("ws://localhost:4000", nil)

	// assertion 1 (check websocket upgrade connection status)
	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %v, want: %v", got, want)
	}

	// register handler for AddPost message
	testPostRouter.RegisterHandler("post add", AddPost)

	// testPostRouter.RegisterHandler("post add", AddPost)

	// creating test message and passing it through websocket
	rawMessage := json.RawMessage(`{"name":"post add", ` +
		`"data":{"name":"Jon", "time":"2017-12-08T12:09:57.341Z", "text":"Hey Neigh!", "address_id":"123hhsj111"}}`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}
	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("posts").Nth(0).Run(session)

	var row map[string]string
	var expectedStr string
	// res.One(&row) <- try and use this thing
	for res.Next(&row) {
		expectedStr = row["name"]
	}
	got2, want2 := expectedStr, "Jon"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	expectedStr = row["time"]

	got2, want2 = expectedStr, "2017-12-08T12:09:57.341Z"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	expectedStr = row["text"]

	got2, want2 = expectedStr, "Hey Neigh!"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	expectedStr = row["address_id"]

	got2, want2 = expectedStr, "123hhsj111"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	r.TableDrop("posts").Wait().Exec(session)

}
