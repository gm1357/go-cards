# go-cards

A small HTTP API for creating, persisting, and manipulating decks of cards. Built as a learning project using only the Go standard library.

## Layout

```
go-cards/
  main.go         entrypoint — starts the HTTP server on :8080
  deck/           card and deck domain logic
    deck.go
    deck_test.go
  server/         HTTP handlers and routing
    server.go
    server_test.go
  decks/          runtime directory where deck files are persisted (created on first POST /deck)
```

Module path: `cards`. Requires Go 1.22+ (uses `http.ServeMux` method-based pattern routing).

## Running

```sh
go run .
```

The server listens on `http://localhost:8080` and persists decks to a `decks/` directory in the working directory.

## Testing

```sh
go test ./...          # run all packages
go test ./deck         # run a single package
go test ./server -v    # verbose output
go test ./... -count=1 # bypass the test cache
```

Server tests spin up an isolated `httptest.Server` with a per-test temp directory, so no shared port and no leftover files.

## API

| Method | Path                          | Description                                                 |
| ------ | ----------------------------- | ----------------------------------------------------------- |
| GET    | `/card/random`                | Return a random card from a fresh, unsaved deck.            |
| POST   | `/deck`                       | Create and persist a new 52-card deck. Returns `{ "id" }`.  |
| GET    | `/deck/{id}`                  | Return the stored deck as JSON.                             |
| GET    | `/deck/{id}/card/random`      | Return a random card from the stored deck.                  |
| GET    | `/deck/{id}/card/top`         | Return the top card of the stored deck.                     |
| POST   | `/deck/{id}/shuffle`          | Shuffle the stored deck in place and return it.             |
| POST   | `/deck/{id}/deal?handSize=N`  | Deal `N` cards off the top. Stored deck shrinks by `N`.     |

### Example

```sh
curl -X POST http://localhost:8080/deck
# {"id":"a1b2c3..."}

curl http://localhost:8080/deck/a1b2c3.../card/top
# {"suit":"Spades","value":"Ace"}

curl -X POST "http://localhost:8080/deck/a1b2c3.../deal?handSize=5"
# [{"suit":"Spades","value":"Ace"}, ...]
```

## Storage

Decks are saved as plain-text files at `decks/<id>.deck`, one card per comma-separated entry in the form `Value:Suit`. The `deck` package exposes `SaveToFile` and `NewFromFile` for the format.
