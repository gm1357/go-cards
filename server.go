package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleCardRandom(w http.ResponseWriter, r *http.Request) {
	d := newDeck()
	c := d.getRandomCard()

	j, err := json.Marshal(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func runServer() {
	http.HandleFunc("/card/random", handleCardRandom)

    log.Fatal(http.ListenAndServe(":8080", nil))
}