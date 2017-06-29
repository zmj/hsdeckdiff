package main

type deck struct {
	Name   string
	Format string
	Cards  []cardCount
	Code   string
}

type card struct {
	Name string
	Cost string
}

type cardCount struct {
	Card  card
	Count int
}

func parse() (deck, bool) {

}
