// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package tmln

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Deck struct {
	Start      *Card
	Paths      []*Card
	Finish     *Card
	CardWidth  int
	CardHeight int
}

// LoadDeck reads the deck configuration from a JSON file and returns a Deck of cards.
// The configuration assigns image files to each card. For Path cards, it also assigns
// the timelines and waypoints to slots 0 through 3 on the card.
func LoadDeck(path string, images string) (*Deck, error) {
	// store defines the structure of the json data
	var store struct {
		Start struct {
			Image string `json:"image"`
		} `json:"start"`
		Finish struct {
			Image string `json:"image"`
		} `json:"finish"`
		Paths []struct {
			Image     string    `json:"image"`
			Lines     [4]int    `json:"lines"`
			Waypoints [4]string `json:"waypoints"`
		} `json:"paths"`
	}

	if data, err := os.ReadFile(path); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	// load the cards included with the distribution
	var d Deck
	d.CardWidth, d.CardHeight = 750, 1049
	d.Start = &Card{Kind: CKStart, Image: filepath.Join(images, store.Start.Image)}
	d.Finish = &Card{Kind: CKFinish, Image: filepath.Join(images, store.Finish.Image)}
	for _, p := range store.Paths {
		c := &Card{Kind: CKPath, Image: filepath.Join(images, p.Image)}
		for slot := 0; slot < 4; slot++ {
			c.Slots[slot] = p.Lines[slot]
			w := &Waypoint{}
			switch p.Waypoints[slot] {
			case "":
				w.Commodity = None
			case "B":
				w.Commodity = OilPetrol
			case "G":
				w.Commodity = Nucleons
			case "R":
				w.Commodity = BeetCandy
			case "Y":
				w.Commodity = MilkBread
			}
			c.Waypoints[slot] = w
		}
		d.Paths = append(d.Paths, c)
	}

	return &d, nil
}
