package db

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
)

type opts struct {
	session *r.Session
	db      string
	table   string
	pK      string
	index   string
}

func connect() *r.Session {
	session, _ := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	return session
}

// createDB Creates a db with the given name
// It will exit if an error occurs
func createDB(opts opts) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	r.DBCreate(opts.db).RunWrite(opts.session)
}

// CreateTable creates the tables feeds, posts, and users
// required for the application to run
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

// TablesDrop drops the tables feeds, posts and users
func dropTable(opts opts) {
	_, err := r.DB(opts.db).TableDrop(opts.table).RunWrite(opts.session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

//DbDrop drops the give database and handles error exiting if there are any
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
	r.DB(opts.db).Table(opts.table).IndexCreate(opts.index).RunWrite(opts.session)
	r.DB(opts.db).Table(opts.table).IndexWait().Run(opts.session)
}

func Setup(dbName string) {
	session := connect()

	createDB(opts{session: session, db: dbName})

	createTable(opts{session: session, db: dbName, table: "users", pK: "email"})
	createTable(opts{session: session, db: dbName, table: "addresses", pK: "postcode"})
	createTable(opts{session: session, db: dbName, table: "feedAddresses", pK: "id"})
	createTable(opts{session: session, db: dbName, table: "feeds", pK: "id"})
	createTable(opts{session: session, db: dbName, table: "posts", pK: "id"})

	createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "feed"})
	createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "address"})

	session.Close()
}

func CleanUp(dbName string) {
	session := connect()
	dropDB(opts{session: session, db: dbName})
	session.Close()
}
