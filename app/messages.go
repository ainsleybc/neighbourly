package app

import "time"

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type User struct {
	Username    string `gorethink:"username,omitempty"`
	Email       string `gorethink:"email,omitempty"`
	Postcode    string `gorethink:"postcode,omitempty"`
	Password    string `gorethink:"password,omitempty"`
	DefaultFeed string `gorethink:"defaultFeed,omitempty"`
}

type Address struct {
	Postcode string `gorethink:"postcode,omitempty"`
}

type FeedAddress struct {
	ID      string `gorethink:"id,omitempty"`
	Address string `gorethink:"address,omitempty"`
	Feed    string `gorethink:"feed,omitempty"`
}

type Feed struct {
	ID      string `gorethink:"id,omitempty"`
	Address string `gorethink:"address,omitempty"`
}

type Post struct {
	ID        string    `gorethink:"id,omitempty"`
	Name      string    `gorethink:"name"`
	CreatedAt time.Time `gorethink:"createdAt"`
	Text      string    `gorethink:"text"`
	FeedID    string    `gorethink:"feedId"`
}
