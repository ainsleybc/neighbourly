package db_setup

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

const DBNAME = "neighbourly"

func init() {
	session, _ = r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	dbCreate(DBNAME)
	tablesCreate(DBNAME)
}

func dbCreate(dbname string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBCreate(dbname).RunWrite(session)
	fmt.Printf("%s db created\n", dbname)
}

func tablesCreate(dbname string) {
	tables := [3]string{"feeds", "posts", "users"}
	for i := 0; i < len(tables); i++ {
		_, err := r.DB(dbname).TableCreate(tables[i]).RunWrite(session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s table created\n", tables[i])
	}
}
