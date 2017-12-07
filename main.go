package main

import (
	"net/http"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

func main() {

	session, _ = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "neighbourly",
	})

	router := NewRouter()

	router.RegisterHandler("feed add", addFeed)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)

}
