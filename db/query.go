/*

query.go is used for all database related queries required by the app.
Exported functions:

TODO

*/

package db

import r "github.com/dancannon/gorethink"

const (
	addresses     = "addresses"
	feeds         = "feeds"
	feedAddresses = "feedAddresses"
	users         = "users"
)

func InsertAddress(session *r.Session, row interface{}) (r.WriteResponse, error) {
	resp, err := r.Table(addresses).
		Insert(row).
		RunWrite(session)
	return resp, err
}

func InsertFeed(session *r.Session, row interface{}) (r.WriteResponse, error) {
	resp, err := r.Table(feeds).
		Insert(row).
		RunWrite(session)
	return resp, err
}

func InsertFeedAddress(session *r.Session, row interface{}) (r.WriteResponse, error) {
	resp, err := r.Table(feedAddresses).
		Insert(row).
		RunWrite(session)
	return resp, err
}

func InsertUser(session *r.Session, row interface{}) (r.WriteResponse, error) {
	resp, err := r.Table(users).
		Insert(row).
		RunWrite(session)
	return resp, err
}

func GetDefaultFeedByAddress(session *r.Session, key interface{}) (*r.Cursor, error) {
	cursor, err := r.Table(feedAddresses).
		GetAllByIndex("address", key).
		EqJoin("feed", r.Table(feeds)).Zip().
		Filter(map[string]interface{}{
			"addressDefault": true,
		}).
		Run(session)
	return cursor, err
}
