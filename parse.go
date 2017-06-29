package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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
	Cost int
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

	type assignFunc func(match []string) error
	parseLine := func(lineName string, reg *regexp.Regexp, assign assignFunc) error {
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
		err := assign(match)
		if err != nil {
			return fmt.Errorf("Failed to assign %v line values: %v", lineName, err.Error())
		}
		return nil
	}

	if err := parseLine("name", p.Name, func(match []string) error {
		d.Name = match[1]
		return nil
	}); err != nil {
		return nil, err
	}

	if err := parseLine("class", p.Class, func(match []string) error {
		d.Class = match[1]
		return nil
	}); err != nil {
		return nil, err
	}

	if err := parseLine("format", p.Format, func(match []string) error {
		d.Format = match[1]
		return nil
	}); err != nil {
		return nil, err
	}

	if err := parseLine("empty", nil, nil); err != nil {
		return nil, err
	}

	deckCardCount := 0
	for {
		if err := parseLine("card", p.CardCount, func(match []string) error {
			cost, err := strconv.Atoi(match[2])
			if err != nil {
				return fmt.Errorf("Failed to parse cost value: %v %v", match[2], err.Error())
			}
			count, err := strconv.Atoi(match[1])
			if err != nil {
				return fmt.Errorf("Failed to parse count value: %v %v", match[1], err.Error())
			}
			card := cardCount{
				Card: card{
					Name: match[3],
					Cost: cost,
				},
				Count: count,
			}
			d.Cards = append(d.Cards, card)
			deckCardCount += card.Count
			return nil
		}); err != nil {
			if strings.TrimSpace(scanner.Text()) != "#" {
				return nil, fmt.Errorf("Unexpected card line: %v %v", scanner.Text(), err.Error())
			}
			if deckCardCount != 30 {
				return nil, fmt.Errorf("Deck does not contain 30 cards")
			}
			break
		}
	}

	if err := parseLine("code", p.Code, func(match []string) error {
		d.Code = match[1]
		return nil
	}); err != nil {
		return nil, err
	}

	return d, nil
}

type parser struct {
	Name      *regexp.Regexp
	Class     *regexp.Regexp
	Format    *regexp.Regexp
	CardCount *regexp.Regexp
	Code      *regexp.Regexp
}

const (
	nameReg   = `### (.+)`
	classReg  = `# Class: (.+)`
	formatReg = `# Format: (.+)`
	cardReg   = `# (\d)x \((\d+)\) (.+)`
	codeReg   = `(.+)`
)

func newParser() (*parser, error) {
	p := &parser{}
	compile := func(regStr, name string, assign func(*regexp.Regexp)) error {
		reg, err := regexp.Compile(regStr)
		if err != nil {
			return fmt.Errorf("Failed to compile %v regexp: %v", name, err.Error())
		}
		assign(reg)
		return nil
	}

	if err := compile(nameReg, "name", func(reg *regexp.Regexp) {
		p.Name = reg
	}); err != nil {
		return nil, err
	}

	if err := compile(classReg, "class", func(reg *regexp.Regexp) {
		p.Class = reg
	}); err != nil {
		return nil, err
	}

	if err := compile(formatReg, "format", func(reg *regexp.Regexp) {
		p.Format = reg
	}); err != nil {
		return nil, err
	}

	if err := compile(cardReg, "card", func(reg *regexp.Regexp) {
		p.CardCount = reg
	}); err != nil {
		return nil, err
	}

	if err := compile(codeReg, "code", func(reg *regexp.Regexp) {
		p.Code = reg
	}); err != nil {
		return nil, err
	}

	return p, nil
}
