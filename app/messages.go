package app

import "time"

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
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
