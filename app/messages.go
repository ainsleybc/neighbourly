package app

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Feed struct {
	Id      string `gorethink:"id,omitempty"`
	Address string `gorethink:"address"`
}
