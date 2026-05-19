package deck

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

type Deck []Card

func New() Deck {
	cards := Deck{}

	suits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	values := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "King", "Queen", "Jack"}

	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, Card{Suit: suit, Value: value})
		}
	}

	return cards
}

func (c Card) String() string {
	return fmt.Sprintf("%v of %v", c.Value, c.Suit)
}

func (d Deck) Print() {
	for i, c := range d {
		fmt.Println(i, c)
	}
}

func (d *Deck) Deal(handSize int) Deck {
	hand := (*d)[:handSize]
	*d = (*d)[handSize:]
	return hand
}

func (d Deck) toString() string {
	var parts []string
	for _, v := range d {
		parts = append(parts, fmt.Sprintf("%v:%v", v.Value, v.Suit))
	}
	return strings.Join(parts, ",")
}

func (d Deck) SaveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0666)
}

func newFromStringArr(s []string) Deck {
	d := Deck{}
	for _, v := range s {
		parts := strings.Split(v, ":")
		d = append(d, Card{Value: parts[0], Suit: parts[1]})
	}
	return d
}

func NewFromFile(filename string) (Deck, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return newFromStringArr(strings.Split(string(bs), ",")), nil
}

func randomPosition(d Deck) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(len(d) - 1)
}

func (d Deck) Shuffle() {
	for i := range d {
		newPos := randomPosition(d)
		d[i], d[newPos] = d[newPos], d[i]
	}
}

func (d Deck) GetRandomCard() Card {
	return d[randomPosition(d)]
}
