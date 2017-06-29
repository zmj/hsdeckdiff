package main

type deckDiff struct {
	from    deck
	to      deck
	Added   []cardCount
	Removed []cardCount
}

func diff(from, to deck) deckDiff {
	return deckDiff{
		from:    from,
		to:      to,
		Added:   nil,
		Removed: nil,
	}
}

func (d *deckDiff) String() string {
	return ""
}
