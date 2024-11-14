// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/playbymail/tmln"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	log.Printf("loading deck\n")
	deck, err := tmln.LoadDeck("../cfl/cards.json", "../cfl/images")
	if err != nil {
		log.Fatalf("error loading deck: %v\n", err)
	}
	log.Printf("deck: %d path cards\n", len(deck.Paths))

	// Verify labels by saving each card as a separate file
	for _, card := range deck.Paths {
		cardPath := fmt.Sprintf("card_debug.png")
		f, err := os.Create(cardPath)
		if err != nil {
			log.Fatalf("error creating card image file: %v", err)
		}
		defer f.Close()
		if err := png.Encode(f, card.Image); err != nil {
			log.Fatalf("error encoding png for card: %v", err)
		}
		break
	}

	// print out a card for debugging
	if card := deck.Paths[0]; card != nil {
		// save the image to a file.
		debugCard := filepath.Join("../cfl", "images", "debug_card.png")
		if f, err := os.Create(debugCard); err != nil {
			log.Fatalf("error creating file: %v\n", err)
		} else if err := png.Encode(f, card.Image); err != nil {
			log.Fatalf("error encoding png: %v\n", err)
		} else {
			_ = f.Close()
		}
		log.Printf("created %s\n", debugCard)
	}

	// print out the deck for debugging
	for page := 0; page < 3; page++ {
		var board [3][3]*tmln.Card
		pageOffset := page * 9
		for row := 0; row < 3; row++ {
			rowOffset := row * 3
			for col := 0; col < 3; col++ {
				offset := pageOffset + rowOffset + col
				if offset >= len(deck.Paths) {
					continue
				}
				board[row][col] = deck.Paths[offset]
			}
		}
		// create an image that is 3 cards wide and 3 cards high.
		img := image.NewRGBA(image.Rect(0, 0, 3*deck.CardWidth, 3*deck.CardHeight))
		// draw the cards on the board.
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				card := board[row][col]
				if card == nil {
					continue
				}
				// compute the position to draw the card on the board
				x, y := col*deck.CardWidth, row*deck.CardHeight
				// draw the card on the board.
				draw.Draw(img, image.Rect(x, y, x+deck.CardWidth, y+deck.CardHeight), card.Image, image.Point{}, draw.Over)
			}
		}
		// save the image to a file.
		debugPage := filepath.Join("../cfl", "images", fmt.Sprintf("debug_page_%02d.png", page+1))
		if f, err := os.Create(debugPage); err != nil {
			log.Fatalf("error creating file: %v\n", err)
		} else if err := png.Encode(f, img); err != nil {
			log.Fatalf("error encoding png: %v\n", err)
		} else {
			_ = f.Close()
		}
		log.Printf("created %s\n", debugPage)
	}
}
