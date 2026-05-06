package main

import "fmt"

func main() {
	cards := newDeck()

	hand, remaining := cards.deal(5)

	fmt.Println("Hand")
	hand.print()
	fmt.Println("Remaining cards")
	remaining.print()
}