# Web Socket Integration points

## All Messages

```json
{
  name: "<message name>",
  data: {
    <message data>
  }
}

```

## data structures

### Feeds

#### client -> server

```json
{
  name: "feed add",
  data: {
    Address: "name"
  }
}
```

#### server -> client

```json
{
  name: "feed add",
  data: {
    ID: "<feed id>"
    Address: "<feed name>"
  }
}
```

### Posts

#### client -> server

```json
{
  name: "add post",
  data: {
    Name: "<post name>"
    Text: "<post text>"
    FeedID: "<feed id>"
  }
}

```

#### server -> client

```json
{
  name: "add post",
  data: {
    ID: "<post id>"
    CreatedAt: "<post createAt timestamp>"
    Name: "<post name>"
    Text: "<post text>"
    FeedID: "<feed id>"
  }
}

```
