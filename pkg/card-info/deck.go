// Copyright 2016-2020 misatos.angel@gmail.com.  All rights reserved.

// Simple package for reading soku card information and decoding them from
// e.g. replay or net-stream game information.

// Use of this source code is governed by GPL-style
// license that can be found in the LICENSE file.

package cardinfo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

type Card struct {
	Code uint16
	Name string
	Type string
	Char string
	Cost uint16
}

type CardCount struct {
	Card  *Card
	Count uint16
}

type Deck struct {
	Total uint16
	Cards map[uint16]*CardCount
}

// Returns an English pretty-name for the given card.
func (self *Card) String() string {
	var usage string
	switch self.Type {
	case "Skill":
		usage = fmt.Sprintf("%d Skill", self.Cost)
	case "System":
		if self.Cost == 0 {
			usage = fmt.Sprintf("Infinite use System")
		} else {
			usage = fmt.Sprintf("Single use System")
		}
	case "Spell":
		usage = fmt.Sprintf("%d-Cost Spell", self.Cost)
	default:
		usage = fmt.Sprintf("[huh? %s]", self.Type)
	}

	return fmt.Sprintf("[%02x] %s Card: '%s'", self.Code, usage, self.Name)
}

// Returns an English pretty-name for the given deck, with each English
// name one per line.
func (self *Deck) String() string {
	keys := make([]int, len(self.Cards))
	i := 0
	for k := range self.Cards {
		keys[i] = int(k)
		i++
	}
	sort.Ints(keys)
	str := ""
	for _, i := range keys {
		cc := self.Cards[uint16(i)]
		count := cc.Count
		card := cc.Card
		str += fmt.Sprintf("%d x %s\n", count, card.String())
	}
	return str
}

// Adds a given card, possibly overtwriting a conflicting card
func (self *Deck) AddCard(c *Card) {
	self.Total++
	existing := self.Cards[c.Code]
	if existing == nil {
		self.Cards[c.Code] = &CardCount{c, 1}
		return
	}
	if existing.Card.Name != c.Name {
		self.Total -= existing.Count
		self.Cards[c.Code] = &CardCount{c, 1}
		return
	}
	existing.Count += 1
}

// returns the count of the number of given cards in the deck
func (self *Deck) GetCardCount(code uint16) uint16 {
	cc := self.Cards[code]
	if cc == nil {
		return 0
	}
	return cc.Count
}

func (self *Deck) GetCard(code uint16) *Card {
	cc := self.Cards[code]
	if cc == nil {
		return nil
	}
	return cc.Card
}

// Attempts to read a given card from the passed reader
// in CSV format. The expected format is:
// Char,Type,Hex,Name,Cost/Input
// Char - Character name or 'System'
// Type - Card type, e.g. 'System', 'Skill' or 'Spell'
// Hex - The hexidecimal (but no '0x' prefix) code for the card
// Name - The pretty name of the card
// Cost/Input - Either the cost in cards (for spell/System) or for skill cards, the actual decimal input, e.g. 236, 214, 22, 623, 421
func NewCardFromCSV(r *csv.Reader) (*Card, error) {
	vals, err := r.Read()
	if err != nil {
		return nil, err
	}
	code, err := strconv.ParseUint(vals[2], 16, 16)
	if err != nil {
		str := fmt.Sprintf("Bad card code value: '%s' in record: '%v'\n", vals[2], vals)
		return nil, errors.New(str)
	}
	cost, err := strconv.ParseUint(vals[4], 10, 16)
	if err != nil {
		str := fmt.Sprintf("Bad card cost value: '%s' in record: '%v'\n", vals[4], vals)
		return nil, errors.New(str)
	}
	return &Card{
		Code: uint16(code),
		Name: vals[3],
		Type: vals[1],
		Char: vals[0],
		Cost: uint16(cost),
	}, nil
}
