package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/posener/wstest"
)

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}

func TestAddFeed(t *testing.T) {

	// connect to rethinkDB
	testDBSession, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "test",
	})

	// close session on end test
	defer testDBSession.Close()

	// create the tables for test
	r.TableCreate("feed").RunWrite(testDBSession)

	// new router
	testRouter := NewRouter(testDBSession)
	d := wstest.NewDialer(testRouter, nil)
	fmt.Println("david")

	// open websocket connection (skip error)
	conn, resp, err := d.Dial("ws://localhost:4000", nil)
	fmt.Println("david")

	// assertion 1 (check websocket upgrade connection status)
	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %v, want: %v", got, want)
	}

	// creating test message and passing it through websocket
	fmt.Println("david")

	testRouter.Handle("feed add", addFeed)
	fmt.Println("david")

	rawMessage := []byte(`{"name":"feed add", ` +
		`"data":{"Address":"Makers Academy"`)

	err = conn.WriteJSON(rawMessage)
	if err != nil {
		t.Fatal(err)
	}

	// simple timeout to allow to database writes
	fmt.Println("david")
	time.Sleep(time.Second * 1)

	// write assertion

	res, err := r.Table("feed").Nth(0).Run(testDBSession)
	// fmt.Printf("%v", conn)
	// fmt.Printf("%v", rawMessage)

	var row map[string]string
	var david string
	fmt.Println("david")

	for res.Next(&row) {
		fmt.Println("david")
		david = row["Address"]
	}
	got2, want2 := david, "Makers Academy"
	if got2 != want2 {
		t.Errorf("got: %v, want: %v", got2, want2)
	}
	r.TableDrop("feed").Wait().Exec(testDBSession)
}

// read assertion needed
