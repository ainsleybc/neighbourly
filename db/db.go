package db

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

func connect() {
	session, _ = r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
}

func close() {
	session.Close()
}

// createDB Creates a db with the given name
// It will exit if an error occurs
func createDB(dbName string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBCreate(dbName).RunWrite(session)
}

// CreateTable creates the tables feeds, posts, and users
// required for the application to run
func createTable(dbName string, tableName string) {
	_, err := r.DB(dbName).TableCreate(tableName).RunWrite(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

// TablesDrop drops the tables feeds, posts and users
func dropTable(dbName string, tableName string) {
	_, err := r.DB(dbName).TableDrop(tableName).RunWrite(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

//DbDrop drops the give database and handles error exiting if there are any
func dropDB(dbName string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBDrop(dbName).RunWrite(session)
}

func Setup(dbName string) {
	connect()
	createDB(dbName)
	createTable(dbName, "users")
	createTable(dbName, "feeds")
	createTable(dbName, "posts")
	close()
}

func CleanUp(dbName string) {
	connect()
	dropDB(dbName)
	close()
}
