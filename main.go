package main

import "fmt"

func main() {
	cards := newDeck()

	hand := cards.deal(5)

	fmt.Println("Hand")
	hand.print()
	fmt.Println("Remaining cards")
	cards.print()
}