package main

import (
	"net/http"

	app "github.com/ainsleybc/neighbourly/app"
	r "github.com/dancannon/gorethink"
)

func main() {

	session, _ := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "neighbourly",
	})

	router := app.NewRouter(session)

	router.RegisterHandler("feed add", app.AddFeed)
	router.RegisterHandler("post add", app.AddPost)

	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)

}
