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
	"io"
)

type AllCards map[string]*Deck

func NewFromCSV(r io.Reader) (AllCards, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = 5
	cards := make(map[string]*Deck, 21)
	first := true
	parsedOne := false
	var firstErr error
	for {
		c, err := NewCardFromCSV(csvReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			if first {
				first = false
				firstErr = err
				continue
			}
			return nil, err
		}
		first = false
		if cards[c.Char] == nil {
			cards[c.Char] = &Deck{0, make(map[uint16]*CardCount, 20)}
		}
		cards[c.Char].AddCard(c)
		parsedOne = true
	}
	if !parsedOne {
		if firstErr != nil {
			return nil, firstErr
		}
		return nil, errors.New("Passed an empty input stream to form a deck from\n")
	}
	return cards, nil
}

func (self *AllCards) FindCard(char string, code uint16) *Card {
	deck := (*self)[char]
	if deck == nil {
		return nil
	}
	return deck.GetCard(code)
}

func (self *AllCards) NewDeck(char string, cards []uint16) (*Deck, error) {
	charDeck := self.FindDeck(char)
	if charDeck == nil {
		return nil, errors.New("Unknown deck: '" + char + "'\n")
	}
	sysDeck := self.FindDeck("System")
	if sysDeck == nil {
		return nil, errors.New("Unable to find System deck\n")
	}

	d := &Deck{0, make(map[uint16]*CardCount, len(cards))}
	for _, code := range cards {
		if code < 100 {
			c := sysDeck.GetCard(code)
			if c == nil {
				str := fmt.Sprintf("Unknown system card id: '%d'\n", code)
				return nil, errors.New(str)
			}
			d.AddCard(c)
			continue
		}
		c := charDeck.GetCard(code)
		if c == nil {
			str := fmt.Sprintf("Unknown '%s' card id: '%d'\n", char, code)
			return nil, errors.New(str)
		}
		d.AddCard(c)
	}
	return d, nil
}

func (self *AllCards) FindDeck(char string) *Deck {
	deck := (*self)[char]
	if deck == nil {
		return nil
	}
	return deck
}
