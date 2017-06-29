package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type deck struct {
	Name   string
	Class  string
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

func (p *parser) parse(deckRaw io.Reader) (*deck, error) {
	scanner := bufio.NewScanner(deckRaw)
	d := &deck{}

	nextLine := func(lineName string) error {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("Failed to read expected %v line: %v", lineName, err.Error())
			}
			return fmt.Errorf("Missing expected %v line", lineName)
		}
		return nil
	}

	parseLine := func(lineName string, reg *regexp.Regexp, assign func([]string)) error {
		if err := nextLine(lineName); err != nil {
			return err
		}
		if reg == nil {
			return nil // expected empty line
		}
		match := reg.FindStringSubmatch(scanner.Text())
		if match == nil || len(match) < reg.NumSubexp()+1 {
			return fmt.Errorf("Failed to parse %v line: %v", lineName, scanner.Text())
		}
		assign(match)
		return nil
	}

	if err := parseLine("name", p.Name, func(match []string) {
		d.Name = match[1]
	}); err != nil {
		return nil, err
	}

	if err := parseLine("class", p.Class, func(match []string) {
		d.Class = match[1]
	}); err != nil {
		return nil, err
	}

	if err := parseLine("format", p.Format, func(match []string) {
		d.Format = match[1]
	}); err != nil {
		return nil, err
	}

	if err := parseLine("empty", nil, nil); err != nil {
		return nil, err
	}

}

type parser struct {
	Name      *regexp.Regexp
	Class     *regexp.Regexp
	Format    *regexp.Regexp
	CardCount *regexp.Regexp
	Code      *regexp.Regexp
}
