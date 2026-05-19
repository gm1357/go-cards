package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var decksDir = "decks"

func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func deckPath(id string) string {
	return filepath.Join(decksDir, id + ".deck")
}

func loadDeck(id string) (deck, error) {
	bs, err := os.ReadFile(deckPath(id))
	if err != nil {
		return nil, err
	}
	return newDeckFromStringArr(strings.Split(string(bs), ",")), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	j, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func handleCardRandom(w http.ResponseWriter, r *http.Request) {
	d := newDeck()
	writeJSON(w, http.StatusOK, d.getRandomCard())
}

func handleDeckCreate(w http.ResponseWriter, r *http.Request) {
	if err := os.MkdirAll(decksDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := newID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := newDeck()
	if err := d.saveToFile(deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func handleDeckGet(w http.ResponseWriter, r *http.Request) {
	d, err := loadDeck(r.PathValue("id"))
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func handleDeckCardRandom(w http.ResponseWriter, r *http.Request) {
	d, err := loadDeck(r.PathValue("id"))
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d.getRandomCard())
}

func handleDeckCardTop(w http.ResponseWriter, r *http.Request) {
	d, err := loadDeck(r.PathValue("id"))
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}
	if len(d) == 0 {
		http.Error(w, "deck is empty", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, d[0])
}

func handleDeckShuffle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := loadDeck(id)
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}

	d.shuffle()
	if err := d.saveToFile(deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func handleDeckDeal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := loadDeck(id)
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}

	handSizeStr := r.URL.Query().Get("handSize")
	handSize, err := strconv.Atoi(handSizeStr)
	if err != nil || handSize <= 0 {
		http.Error(w, "invalid handSize", http.StatusBadRequest)
		return
	}
	if handSize > len(d) {
		http.Error(w, "handSize exceeds deck size", http.StatusBadRequest)
		return
	}

	hand := d.deal(handSize)
	if err := d.saveToFile(deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, hand)
}

func runServer() {
	http.HandleFunc("GET /card/random", handleCardRandom)
	http.HandleFunc("POST /deck", handleDeckCreate)
	http.HandleFunc("GET /deck/{id}", handleDeckGet)
	http.HandleFunc("GET /deck/{id}/card/random", handleDeckCardRandom)
	http.HandleFunc("GET /deck/{id}/card/top", handleDeckCardTop)
	http.HandleFunc("POST /deck/{id}/shuffle", handleDeckShuffle)
	http.HandleFunc("POST /deck/{id}/deal", handleDeckDeal)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
