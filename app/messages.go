package app

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Feed struct {
	Id      string `gorethink:"id,omitempty"`
	Address string `gorethink:"address"`
}

type Post struct {
	Id         string `gorethink:"id,omitempty"`
	Name       string `gorethink:"name"`
	Time       string `gorethink:"time"`
	Text       string `gorethink:"text"`
	Address_id string `gorethink:"address_id"`
}
