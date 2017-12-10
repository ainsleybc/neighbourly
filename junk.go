package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
)

func main() {

	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})

	fmt.Println(session)

	err := r.DB("neighbourly").Exec(session)

	fmt.Println(err)
}
