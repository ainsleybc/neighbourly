package main

import (
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
	testDBSession, err := r.Connect(r.ConnectOps{
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

	// open websocket connection (skip error)
	conn, resp, _ := d.Dial("ws://localhost:4000", nil)

	// assertion 1 (check websocket upgrade connection status)
	got, want := resp.StatusCode, http.StatusSwitchingProtocols
	if got != want {
		t.Errorf("resp.StatusCode: %v, want: %v", got, want)
	}

	// creating test message and passing it through websocket
	testRouter.Handle("feed add", addFeed)
	rawMessage := []byte(`{"name":"feed add", ` +
		`"data":{"Address":"Makers Academy"`)

	err = conn.WriteJson(rawMessage)
	if err != nil {
		t.Fatal(err)
	}

	// simple timeout to allow to database writes
	time.Sleep(time.Second * 1)

	// write assertion
	res, err := r.Table("feed").Run(testDBSession)
	got, want := res["Address"], "Makers Academy"
	if got != want {
		t.Errorf("got: %v, wnat: %v", got, want)
	}
	r.TableDrop("feed").Wait().Exec(testSession)
}

// read assertion needed
