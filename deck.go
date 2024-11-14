// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package tmln

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	_ "image/png" // register PNG decoder
	"log"
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
//
// Path is the JSON configuration file.
// Images is the directory containing the image files.
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

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	gof, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(gof, &truetype.Options{Size: 48})

	// load the cards included with the distribution
	var d Deck
	d.CardWidth, d.CardHeight = 750, 1049
	d.Start = &Card{Kind: CKStart}
	d.Start.Image, err = loadImage(images, store.Start.Image)
	if err != nil {
		return nil, err
	}
	d.Finish = &Card{Kind: CKFinish}
	d.Finish.Image, err = loadImage(images, store.Finish.Image)
	for _, p := range store.Paths {
		c := &Card{Kind: CKPath}
		var waypoints [4]Commodity
		for i, w := range p.Waypoints {
			switch w {
			case "":
				waypoints[i] = None
			case "B":
				waypoints[i] = OilPetrol
			case "G":
				waypoints[i] = Nucleons
			case "R":
				waypoints[i] = BeetCandy
			case "Y":
				waypoints[i] = MilkBread
			}
		}
		c.Image, err = loadImageWithSlots(images, p.Image, p.Lines, waypoints, face)
		if err != nil {
			return nil, err
		}
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

func loadImage(path, file string) (*image.RGBA, error) {
	if data, err := os.ReadFile(filepath.Join(path, file)); err != nil {
		return nil, err
	} else if img, _, err := image.Decode(bytes.NewReader(data)); err != nil {
		return nil, err
	} else {
		// return a new RGBA image
		return image.NewRGBA(img.Bounds()), nil
	}
}

// loadImageWithSlots loads an image from a file, and draws numbers from slots along the right side.
func loadImageWithSlots(path, file string, slots [4]int, waypoints [4]Commodity, face font.Face) (*image.RGBA, error) {
	// Read the image file
	data, err := os.ReadFile(filepath.Join(path, file))
	if err != nil {
		return nil, err
	}

	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Create a gg context from the decoded image
	dc := gg.NewContextForImage(img)
	dc.SetFontFace(face)

	// Set up properties for drawing the text
	dc.SetRGB(1, 1, 1) // White color

	var labels [4]struct {
		LX, RX, Y float64
	}
	for i := range slots {
		labels[i].LX, labels[i].RX = 110, 700
		switch i {
		case 0:
			labels[i].Y = 125
		case 1:
			labels[i].Y = 390
		case 2:
			labels[i].Y = 650
		case 3:
			labels[i].Y = 910
		}
	}
	for i := range slots {
		dc.DrawStringAnchored(fmt.Sprintf("%d%s", i+1, waypoints[i]), labels[i].LX, labels[i].Y, 1, 0.5) // Anchored to the right, center vertically
		dc.DrawStringAnchored(fmt.Sprintf("%d", slots[i]+1), labels[i].RX, labels[i].Y, 1, 0.5)          // Anchored to the right, center vertically
	}

	return dc.Image().(*image.RGBA), nil
}
