package main

import (
	"os"
	"testing"
)

const LEN_DECK = 52
const FIRST_CARD = "Ace of Spades"
const LAST_CARD = "Jack of Clubs"
const TEST_FILE_NAME = "_decktesting"

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != LEN_DECK {
		t.Errorf("Expected deck length of %v, but got %v", LEN_DECK, len(d))
	}

	if d[0].String() != FIRST_CARD {
		t.Errorf("Expected first card to be %v, but got %v", FIRST_CARD, d[0])
	}

	if d[len(d) - 1].String() != LAST_CARD {
		t.Errorf("Expected first card to be %v, but got %v", LAST_CARD, d[len(d) - 1])
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	os.Remove(TEST_FILE_NAME)

	deck := newDeck()
	deck.saveToFile(TEST_FILE_NAME)

	loadedDeck := newDeckFromFile(TEST_FILE_NAME)

	if len(loadedDeck) != len(deck) {
		t.Errorf("Expected deck length of %v, but got %v", len(deck), len(loadedDeck))
	}

	os.Remove(TEST_FILE_NAME)
}