package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

const TEST_DECKS_DIR = "decks_test"
const TEST_SERVER_URL = "http://localhost:8080"

var setupOnce sync.Once

func setup() {
	setupOnce.Do(func() {
		decksDir = TEST_DECKS_DIR
		os.RemoveAll(decksDir)
		go runServer()
		time.Sleep(100 * time.Millisecond)
	})
}

func createDeck(t *testing.T) string {
	t.Helper()
	res, err := http.Post(TEST_SERVER_URL + "/deck", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to create deck: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, but got %v", res.StatusCode)
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if body.ID == "" {
		t.Fatalf("Expected an id, but got empty string")
	}
	return body.ID
}

func TestHandleCardRandom(t *testing.T) {
	setup()
	res, err := http.Get(TEST_SERVER_URL + "/card/random")
	if err != nil {
		t.Fatalf("Failed to get random card: %v", err)
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

func TestHandleDeckCreate(t *testing.T) {
	setup()
	id := createDeck(t)

	if _, err := os.Stat(deckPath(id)); err != nil {
		t.Errorf("Expected deck file to exist for id %v: %v", id, err)
	}
}

func TestHandleDeckGet(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Get(TEST_SERVER_URL + "/deck/" + id)
	if err != nil {
		t.Fatalf("Failed to get deck: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var d deck
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(d) != 52 {
		t.Errorf("Expected deck of 52 cards, but got %v", len(d))
	}
}

func TestHandleDeckGetNotFound(t *testing.T) {
	setup()
	res, err := http.Get(TEST_SERVER_URL + "/deck/nonexistent")
	if err != nil {
		t.Fatalf("Failed to get deck: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, but got %v", res.StatusCode)
	}
}

func TestHandleDeckCardRandom(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Get(TEST_SERVER_URL + "/deck/" + id + "/card/random")
	if err != nil {
		t.Fatalf("Failed to get random card: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var c card
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if c.Suit == "" || c.Value == "" {
		t.Errorf("Expected a valid card, but got %+v", c)
	}
}

func TestHandleDeckCardTop(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Get(TEST_SERVER_URL + "/deck/" + id + "/card/top")
	if err != nil {
		t.Fatalf("Failed to get top card: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var c card
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if c.Suit != "Spades" || c.Value != "Ace" {
		t.Errorf("Expected Ace of Spades as top card, but got %+v", c)
	}
}

func TestHandleDeckShuffle(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Post(TEST_SERVER_URL + "/deck/" + id + "/shuffle", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to shuffle: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var shuffled deck
	if err := json.NewDecoder(res.Body).Decode(&shuffled); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(shuffled) != 52 {
		t.Errorf("Expected 52 cards after shuffle, but got %v", len(shuffled))
	}

	original := newDeck()
	sameOrder := true
	for i := range original {
		if original[i] != shuffled[i] {
			sameOrder = false
			break
		}
	}
	if sameOrder {
		t.Errorf("Expected deck order to change after shuffle")
	}
}

func TestHandleDeckDeal(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Post(TEST_SERVER_URL + "/deck/" + id + "/deal?handSize=5", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to deal: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %v", res.StatusCode)
	}

	var hand deck
	if err := json.NewDecoder(res.Body).Decode(&hand); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(hand) != 5 {
		t.Errorf("Expected hand of 5 cards, but got %v", len(hand))
	}

	getRes, err := http.Get(TEST_SERVER_URL + "/deck/" + id)
	if err != nil {
		t.Fatalf("Failed to get deck after deal: %v", err)
	}
	defer getRes.Body.Close()

	var d deck
	if err := json.NewDecoder(getRes.Body).Decode(&d); err != nil {
		t.Fatalf("Failed to decode deck: %v", err)
	}
	if len(d) != 47 {
		t.Errorf("Expected 47 cards remaining after dealing 5, but got %v", len(d))
	}
}

func TestHandleDeckDealInvalidHandSize(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Post(TEST_SERVER_URL + "/deck/" + id + "/deal?handSize=abc", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to deal: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, but got %v", res.StatusCode)
	}
}

func TestHandleDeckDealHandSizeExceedsDeck(t *testing.T) {
	setup()
	id := createDeck(t)

	res, err := http.Post(TEST_SERVER_URL + "/deck/" + id + "/deal?handSize=100", "application/json", nil)
	if err != nil {
		t.Fatalf("Failed to deal: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400, but got %v", res.StatusCode)
	}
}
