package main

import (
	"fmt"
	"log"

	"cards/server"
)

func main() {
	s := server.New("decks")
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(s.ListenAndServe(":8080"))
}
