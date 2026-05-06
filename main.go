package main

func main() {
	cards := deck{"Ace of Diamonds", "Five of Diamonds"}
	cards = append(cards, "Six of Spades")

	cards.print()
}