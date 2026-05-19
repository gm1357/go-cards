package main

import (
	"fmt"
	"log"

	_ "cards/docs"
	"cards/server"
)

// @title           go-cards API
// @version         1.0
// @description     A small HTTP API for creating, persisting, and manipulating decks of cards.
// @host            localhost:8080
// @BasePath        /
func main() {
	s := server.New("decks")
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(s.ListenAndServe(":8080"))
}
