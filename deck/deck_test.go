package deck

import (
	"os"
	"testing"
)

const lenDeck = 52
const firstCard = "Ace of Spades"
const lastCard = "Jack of Clubs"
const testFileName = "_decktesting"

func TestNew(t *testing.T) {
	d := New()

	if len(d) != lenDeck {
		t.Errorf("Expected deck length of %v, but got %v", lenDeck, len(d))
	}

	if d[0].String() != firstCard {
		t.Errorf("Expected first card to be %v, but got %v", firstCard, d[0])
	}

	if d[len(d)-1].String() != lastCard {
		t.Errorf("Expected last card to be %v, but got %v", lastCard, d[len(d)-1])
	}
}

func TestSaveToFileAndNewFromFile(t *testing.T) {
	os.Remove(testFileName)

	d := New()
	if err := d.SaveToFile(testFileName); err != nil {
		t.Fatalf("Failed to save deck: %v", err)
	}

	loaded, err := NewFromFile(testFileName)
	if err != nil {
		t.Fatalf("Failed to load deck: %v", err)
	}

	if len(loaded) != len(d) {
		t.Errorf("Expected deck length of %v, but got %v", len(d), len(loaded))
	}

	os.Remove(testFileName)
}
