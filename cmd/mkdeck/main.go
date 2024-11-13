// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/tmln"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	log.Printf("loading deck\n")
	deck, err := tmln.LoadDeck("../cfl/cards.json", "../cfl/images")
	if err != nil {
		log.Fatalf("error loading deck: %v\n", err)
	}
	log.Printf("deck: %d path cards\n", len(deck.Paths))
}
