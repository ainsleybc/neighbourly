# Neighbour.ly: Server

## Overview
This is the back-end of Neighbour.ly, a communications app produced as part of our final two-week project of Makers Academy.

To see the ReactJS front-end, [click here](https://github.com/alexiscarlier/neighbourly-app).

Neighbour.ly was produced as part of a wider challenge to learn how to use Go without prior knowledge. To see a record of our learning process and how we came to build Neighbour.ly, [click here](https://github.com/haletothewood/LearningGoAndReact)

## Authors

- David Halewood
- Alexis Carlier
- George Lamprakis
- Jon Sanders
- Lucas Salmins
- Ainsley Chang

## Instructions

Go must be installed and your workspace configured to use this repo. For instructions on this, [click here](https://golang.org/doc/install).

### Running Tests

```
$ go get github.com/ainsleybc/neighbourly
$ cd src/github.com/ainsleybc/neighbourly
$ go test
```

### Manually Run Server

```
$ go get github.com/ainsleybc/neighbourly
$ cd src/github.com/ainsleybc/neighbourly
$ go run main.go
```

With the server running, you can then manually simulate messages sent from the front-end with JavaScript, using the console in your web browser or a service like [JSBin]("https://jsbin.com").

For example:
```
var ws = new WebSocket("ws://localhost:4000")
ws.send('{"name": "feed add","data": {"address":"Makers Academy"}')
```

## Technologies used

#### Go
Main server-side language

#### RethinkDB
Database

#### External packages
- [GoRethink](https://github.com/GoRethink/gorethink): RethinkDB Driver for Go
- [wstest](https://github.com/posener/wstest): Client for testing WebSocket connections in Go
- [Gorilla websocket](https://github.com/gorilla/websocket): A Websocket implementation for Go
- [Map Structure](https://github.com/mitchellh/mapstructure): A Go library for decoding map values into structs
- [Bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt): Adaptive hasing algorithm for passwords


## File Manifest

```
|-- neighbourly
    |-- README.md
    |-- integration.md
    |-- main.go
    |-- app
    |   |-- client.go
    |   |-- handlers.go
    |   |-- messages.go
    |   |-- router.go
    |-- test
        |-- LoginUser_test.go
        |-- addFeed_test.go
        |-- addPost_test.go
        |-- signUpUser_test.go
        |-- subscribeFeed_test.go
        |-- subscribePost_test.go
        |-- rethinkdb_data
            |-- log_file
            |-- metadata
```
