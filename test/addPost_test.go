package app

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/ainsleybc/neighbourly/app"
	"github.com/ainsleybc/neighbourly/db"
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

	t.Parallel()

	db.CleanUp("addPost")
	db.Setup("addPost")
	defer db.CleanUp("addPost")

	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "addPost",
	})

	// close session on end test
	defer session.Close()

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
		`"data":{"name":"Jon", "text":"Hey Neigh!", "feedId":"123hhsj111"}}`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}
	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("posts").Nth(0).Run(session)

	var row map[string]interface{}
	var expectedStr string
	// res.One(&row) <- try and use this thing

	for res.Next(&row) {
		expectedStr = row["name"].(string)
	}

	got2, want2 := expectedStr, "Jon"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	expectedStr = row["text"].(string)

	got2, want2 = expectedStr, "Hey Neigh!"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

	expectedStr = row["feedId"].(string)

	got2, want2 = expectedStr, "123hhsj111"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}

}
