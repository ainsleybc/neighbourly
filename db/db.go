/*

Package db is used for all database related commands required for setting up environments, testing & deployments
Exported functions:

- Setup(dbName string) // creates a new database & applies the schema
- CleanUp(dbName string) // drops a database

TODO - create a migration script for production use.. no data loss

*/
package db

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
)

// opts are used to specify parameters for all db commands
type opts struct {
	session *r.Session
	db      string
	table   string
	pK      string
	index   string
}

// connect starts a database session
func connect() *r.Session {
	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	return session
}

// createDB creates an empty db
func createDB(opts opts) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBCreate(opts.db).RunWrite(opts.session)
}

// CreateTable creates a new table
func createTable(opts opts) {
	_, err := r.DB(opts.db).
		TableCreate(opts.table, r.TableCreateOpts{
			PrimaryKey: opts.pK,
		}).
		RunWrite(opts.session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

// dropTable drops a table
func dropTable(opts opts) {
	_, err := r.DB(opts.db).TableDrop(opts.table).RunWrite(opts.session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

// dropDB drops the the database
func dropDB(opts opts) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBDrop(opts.db).RunWrite(opts.session)
}

func createIndex(opts opts) {
	r.Table(opts.table).IndexCreate(opts.index).RunWrite(opts.session)
}

// Setup sets up a new database
func Setup(dbName string) {
	session := connect()

	createDB(opts{session: session, db: dbName})

	createTable(opts{session: session, db: dbName, table: "users", pK: "email"})
	createTable(opts{session: session, db: dbName, table: "addresses", pK: "postcode"})
	createTable(opts{session: session, db: dbName, table: "feedAddresses"})
	createTable(opts{session: session, db: dbName, table: "feeds"})
	createTable(opts{session: session, db: dbName, table: "posts"})

	createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "feed"})
	createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "address"})

	session.Close()
}

// CleanUp drops the whole database
func CleanUp(dbName string) {
	session := connect()
	dropDB(opts{session: session, db: dbName})
	session.Close()
}
