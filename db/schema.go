/*

Package db is used for all database related commands required for setting up environments, testing & deployments
Exported functions:

- Setup(dbName string) // creates a new database & applies the schema
- CleanUp(dbName string) // drops a database

TODO - create a migration script for production use.. no data loss

*/
package db

import (
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
func createDB(opts opts) error {
	_, err := r.DBCreate(opts.db).RunWrite(opts.session)
	return err
}

// CreateTable creates a new table
func createTable(opts opts) error {
	_, err := r.DB(opts.db).
		TableCreate(opts.table, r.TableCreateOpts{
			PrimaryKey: opts.pK,
		}).
		RunWrite(opts.session)
	return err
}

// dropTable drops a table
func dropTable(opts opts) error {
	_, err := r.DB(opts.db).TableDrop(opts.table).RunWrite(opts.session)
	return err
}

// dropDB drops the the database
func dropDB(opts opts) error {
	_, err := r.DBDrop(opts.db).RunWrite(opts.session)
	return err
}

func createIndex(opts opts) error {
	_, err := r.DB(opts.db).Table(opts.table).IndexCreate(opts.index).RunWrite(opts.session)
	_, err = r.DB(opts.db).Table(opts.table).IndexWait().Run(opts.session)
	return err
}

// Setup sets up a new database
func Setup(dbName string) error {
	session := connect()
	var err error

	err = createDB(opts{session: session, db: dbName})
	if err != nil {
		return err
	}

	err = createTable(opts{session: session, db: dbName, table: "users", pK: "email"})
	err = createTable(opts{session: session, db: dbName, table: "addresses", pK: "address"})
	err = createTable(opts{session: session, db: dbName, table: "feedAddresses", pK: "id"})
	err = createTable(opts{session: session, db: dbName, table: "feeds", pK: "id"})
	err = createTable(opts{session: session, db: dbName, table: "posts", pK: "id"})

	err = createIndex(opts{session: session, db: dbName, table: "posts", index: "createdAt"})
	err = createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "feed"})
	err = createIndex(opts{session: session, db: dbName, table: "feedAddresses", index: "address"})

	session.Close()
	return err
}

// CleanUp drops the whole database
func CleanUp(dbName string) error {
	session := connect()
	var err error
	err = dropTable(opts{session: session, db: dbName, table: "posts"})
	err = dropTable(opts{session: session, db: dbName, table: "addresses"})
	err = dropTable(opts{session: session, db: dbName, table: "users"})
	err = dropTable(opts{session: session, db: dbName, table: "feeds"})
	err = dropTable(opts{session: session, db: dbName, table: "feedAddresses"})
	err = dropDB(opts{session: session, db: dbName})
	session.Close()
	return err
}
