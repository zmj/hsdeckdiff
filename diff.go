package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type deckDiff struct {
	from    *deck
	to      *deck
	Added   []cardCount
	Removed []cardCount
}

func diff(from, to *deck) deckDiff {
	cards := make(map[card]*[2]int)
	for _, c := range from.Cards {
		cards[c.Card] = &[2]int{c.Count, 0}
	}
	for _, c := range to.Cards {
		d, exists := cards[c.Card]
		if exists {
			d[1] = c.Count
		} else {
			cards[c.Card] = &[2]int{0, c.Count}
		}
	}
	var added diffCardCount
	var removed diffCardCount
	for c, counts := range cards {
		if counts[0] == counts[1] {
			continue
		}
		if counts[0] > counts[1] {
			removed = append(removed, cardCount{c, counts[0] - counts[1]})
		}
		if counts[1] > counts[0] {
			added = append(added, cardCount{c, counts[1] - counts[0]})
		}
	}
	sort.Sort(added)
	sort.Sort(removed)
	return deckDiff{
		from:    from,
		to:      to,
		Added:   added,
		Removed: removed,
	}
}

func (d *deckDiff) String() string {
	var buffer bytes.Buffer
	s := fmt.Sprintf("%v -> %v\n", d.from.Name, d.to.Name)
	buffer.WriteString(s)
	for _, add := range d.Added {
		s := fmt.Sprintf("+%v %v\n", add.Count, add.Card.Name)
		buffer.WriteString(s)
	}
	for _, remove := range d.Removed {
		s := fmt.Sprintf("-%v %v\n", remove.Count, remove.Card.Name)
		buffer.WriteString(s)
	}
	return strings.TrimSpace(buffer.String())
}

type diffCardCount []cardCount

func (cards diffCardCount) Len() int {
	return len(cards)
}

func (cards diffCardCount) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}

func (cards diffCardCount) Less(i, j int) bool {
	c := cards[i]
	d := cards[j]

	if c.Card.Cost != d.Card.Cost {
		return c.Card.Cost < d.Card.Cost
	}

	names := []string{c.Card.Name, d.Card.Name}
	return sort.StringsAreSorted(names)
}
