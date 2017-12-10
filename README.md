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

## Technologies used

#### Go
Main server-side language

#### RethinkDB
Database

#### External packages
- [GoRethink](https://github.com/GoRethink/gorethink): RethinkDB Driver for Go
- [wstest](https://github.com/posener/wstest): Client for testing WebSocket connections in Go



## File Manifest

```
|-- neighbourly
    |-- README.md
    |-- integration.md
    |-- junk.go
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
