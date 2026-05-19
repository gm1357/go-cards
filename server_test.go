package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func setup() {
	go runServer()
	time.Sleep(100 * time.Millisecond)
}

func TestHandleCardRandom(t *testing.T) {
	setup()
	res, err := http.Get("http://localhost:8080/card/random")
	if err != nil {
		t.Fatalf("Failed to get random card: %v", err)

		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var c card
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if c.String() == "" {
		t.Errorf("Expected a valid card, but got an empty string")
	}
}