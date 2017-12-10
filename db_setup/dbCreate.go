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

}

// DbCreate Creates a db with the given name
// It will exit if an error occurs
func DbCreate(dbname string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBCreate(dbname).RunWrite(session)
	fmt.Printf("%s db created\n", dbname)
}

// TablesCreate creates the tables feeds, posts, and users
// required for the application to run
func TablesCreate(dbname string) {
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

// TablesDrop drops the tables feeds, posts and users
func TablesDrop(dbname string) {
	tables := [3]string{"feeds", "posts", "users"}
	for i := 0; i < len(tables); i++ {
		_, err := r.DB(dbname).TableDrop(tables[i]).RunWrite(session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s table created\n", tables[i])
	}
}

//DbDrop drops the give database and handles error exiting if there are any
func DbDrop(dbname string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBDrop(dbname).RunWrite(session)
}
