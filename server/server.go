package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"cards/deck"
)

type Server struct {
	decksDir string
	mux      *http.ServeMux
}

func New(decksDir string) *Server {
	s := &Server{decksDir: decksDir, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /card/random", s.handleCardRandom)
	s.mux.HandleFunc("POST /deck", s.handleDeckCreate)
	s.mux.HandleFunc("GET /deck/{id}", s.handleDeckGet)
	s.mux.HandleFunc("GET /deck/{id}/card/random", s.handleDeckCardRandom)
	s.mux.HandleFunc("GET /deck/{id}/card/top", s.handleDeckCardTop)
	s.mux.HandleFunc("POST /deck/{id}/shuffle", s.handleDeckShuffle)
	s.mux.HandleFunc("POST /deck/{id}/deal", s.handleDeckDeal)
}

func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Server) deckPath(id string) string {
	return filepath.Join(s.decksDir, id + ".deck")
}

func (s *Server) loadDeck(id string) (deck.Deck, error) {
	return deck.NewFromFile(s.deckPath(id))
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

func (s *Server) handleCardRandom(w http.ResponseWriter, r *http.Request) {
	d := deck.New()
	writeJSON(w, http.StatusOK, d.GetRandomCard())
}

func (s *Server) handleDeckCreate(w http.ResponseWriter, r *http.Request) {
	if err := os.MkdirAll(s.decksDir, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := newID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := deck.New()
	if err := d.SaveToFile(s.deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleDeckGet(w http.ResponseWriter, r *http.Request) {
	d, err := s.loadDeck(r.PathValue("id"))
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (s *Server) handleDeckCardRandom(w http.ResponseWriter, r *http.Request) {
	d, err := s.loadDeck(r.PathValue("id"))
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, d.GetRandomCard())
}

func (s *Server) handleDeckCardTop(w http.ResponseWriter, r *http.Request) {
	d, err := s.loadDeck(r.PathValue("id"))
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

func (s *Server) handleDeckShuffle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.loadDeck(id)
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}

	d.Shuffle()
	if err := d.SaveToFile(s.deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (s *Server) handleDeckDeal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.loadDeck(id)
	if err != nil {
		http.Error(w, "deck not found", http.StatusNotFound)
		return
	}

	handSize, err := strconv.Atoi(r.URL.Query().Get("handSize"))
	if err != nil || handSize <= 0 {
		http.Error(w, "invalid handSize", http.StatusBadRequest)
		return
	}
	if handSize > len(d) {
		http.Error(w, "handSize exceeds deck size", http.StatusBadRequest)
		return
	}

	hand := d.Deal(handSize)
	if err := d.SaveToFile(s.deckPath(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, hand)
}
