package main

import db "github.com/ainsleybc/neighbourly/db"

func main() {
	db.Setup("neighbourly")
	db.Setup("test")
}
