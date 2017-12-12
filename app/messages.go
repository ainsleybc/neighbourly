package app

import "time"

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type User struct {
	Email       string  `gorethink:"email,omitempty"`
	Username    string  `gorethink:"username,omitempty"`
	Address     Address `gorethink:"address,reference" gorethink_ref:"address"`
	Password    string  `gorethink:"password,omitempty"`
	DefaultFeed Feed    `gorethink:"defaultFeed,reference" gorethink_ref:"id"`
}

type Address struct {
	StreetNumber string `gorethink:"address[0],omitempty"`
	StreetName   string `gorethink:"address[1],omitempty"`
	Postcode     string `gorethink:"address[2],omitempty"`
}

type FeedAddress struct {
	ID      string  `gorethink:"id,omitempty"`
	Address Address `gorethink:"address,reference" gorethink_ref:"address"`
	Feed    Feed    `gorethink:"feed,reference" gorethink_ref:"id"`
}

type Feed struct {
	ID             string `gorethink:"id,omitempty"`
	Name           string `gorethink:"name,omitempty"`
	AddressDefault bool   `gorethink:"addressDefault"`
}

type Post struct {
	ID        string    `gorethink:"id,omitempty"`
	Name      string    `gorethink:"name"`
	CreatedAt time.Time `gorethink:"createdAt"`
	Text      string    `gorethink:"text"`
	Feed      Feed      `gorethink:"feed,reference" gorethink_ref:"id"`
}
