// Copyright 2016-2020 misatos.angel@gmail.com.  All rights reserved.

// simple testing command line program that will take a character and string of hex values
// and prtty-print the names of the cards

// Use of this source code is governed by GPL-style
// license that can be found in the LICENSE file.


package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/soku-cardinfo/pkg/card-info"
	"log"
	"os"
	"strconv"
	"strings"
)

// Variables used for command line parameters
var settings struct {
	Character string   `short:"c" long:"char"  required:"true" description:"Character which cards belong to"`
	CSVFile   string   `short:"f" long:"file" required:"true" description:"Location of a cards CSV file to read"`
	Hex       bool     `short:"H" long:"hex" description:"Presume inputs card values are hexidecimal"`
	Cards     []uint16 `positional-arg-name:"Cards"`
}

func init() {
}

func main() {
	os.Exit(run())
}

func run() int {
	CliParse()
	if len(settings.Cards) == 0 {
		fmt.Printf("Must pass at least one Card\n")
		return 1
	}
	csvFile, err := os.Open(settings.CSVFile)
	if err != nil {
		log.Fatal("Unable to open CSV file:", err)
	}

	allCards, err := cardinfo.NewFromCSV(csvFile)
	if err != nil {
		log.Fatal("Unable to read CSV file:", err)
	}
	deck, err := allCards.NewDeck(settings.Character, settings.Cards)
	if err != nil {
		log.Fatal("Unable to interpret deck:", err)
	}
	fmt.Printf("Deck for %s:\n%s", settings.Character, deck.String())
	return 0
}

func CliParse() {
	parser := flags.NewParser(&settings, flags.Default)
	args, err := parser.Parse()

	if err != nil {
		switch err.(type) {
		case *flags.Error:
			if err.(*flags.Error).Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		log.Fatalln(err)
	}
	l := len(args)
	if len(args) == 0 {
		log.Fatalln("Must pass a card number to lookup")
	}
	settings.Cards = make([]uint16, l)
	for i, codeStr := range args {
		lc_vers := strings.ToLower(codeStr)
		hex := settings.Hex
		base := 10
		if strings.HasPrefix(lc_vers, "0x") {
			hex = true
			codeStr = codeStr[2:]
		}
		if hex {
			base = 16
		}
		val, err := strconv.ParseUint(codeStr, base, 16)
		if err != nil {
			log.Fatalln("Bad (non-integer) code", codeStr, "passed:", err)
		}
		settings.Cards[i] = uint16(val)
	}
}
