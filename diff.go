package main

type deckDiff struct {
	from    deck
	to      deck
	Added   []cardCount
	Removed []cardCount
}

func diff(from, to deck) deckDiff {

}
