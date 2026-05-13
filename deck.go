package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type card struct {
	suit string
	value string
}

type deck []card


func newDeck() deck {
	cards := deck{}

	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "King", "Queen", "Jack"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			c := card{suit: suit, value: value}
			cards = append(cards, c)
		}
	}

	return cards
}

func (c card) String() string {
	return fmt.Sprintf("%v of %v", c.value, c.suit)
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

func (d *deck) deal(handsize int) (deck) {
	hand := (*d)[:handsize]
	*d = (*d)[handsize:]

	return hand
}

func (d deck) toString() string {
	var deckStrings []string

	for _, v := range d {
		deckStrings = append(deckStrings, fmt.Sprintf("%v:%v", v.value, v.suit))
	}

	return strings.Join(deckStrings, ",")
}

func (d deck) saveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0666)
}

func newDeckFromStringArr(s []string) deck {
	d := deck{}

	for _, v := range s {
		cSplit := strings.Split(v, ":")
		c := card{value: cSplit[0], suit: cSplit[1]}
		d = append(d, c)
	}

	return d
}

func newDeckFromFile(filename string) deck {
	bs, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	s := strings.Split(string(bs), ",")
	return newDeckFromStringArr(s)
}

func (d deck) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPosition := r.Intn(len(d) - 1)

		d[i], d[newPosition] = d[newPosition], d[i]
	}
}