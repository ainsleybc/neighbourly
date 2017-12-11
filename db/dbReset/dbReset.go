package main

import db "github.com/ainsleybc/neighbourly/db"

func main() {
	db.CleanUp("neighbourly")
	db.Setup("neighbourly")
}
