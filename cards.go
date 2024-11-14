// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package tmln

import "image"

// there are 23 path cards in the game. we need 46.
// we rotate the original cards 180 degrees to make the missing 23 cards.

// Card represents a card in the game.
//
// All cards have four TimeLines.
//  1. The links for Start cards have nil values for From.
//  2. The Waypoints for TimeLines on Start cards are always empty.
//  3. The links for Finish cards have nil values for To.
//  4. The Waypoints for TimeLInes on Finish cards are always empty.
//  5. The links for Path cards will usually have values for both From and To.
//     Sometimes those links will skip columns on the board.
//  6. The Waypoints for TimeLines on Path cards are initially populated.
//     Players may collect the Commodity tokens during play.
//
// Cards may be played face up or down.
type Card struct {
	Kind      CardKind
	Slots     [4]int       // 0...3, maps input to output
	Timelines [4]*TimeLine // 0...3
	Waypoints [4]*Waypoint // 0...3
	Image     *image.RGBA  // stores the actual decoded PNG
	Rotate180 bool         // when true, rotate the image 180 degrees
}

type CardKind int

const (
	CKStart CardKind = iota
	CKPath
	CKFinish
)

// Commodity represents a commodity token.
type Commodity int

const (
	None Commodity = iota
	BeetCandy
	MilkBread
	Nucleons
	OilPetrol
)

func (c Commodity) String() string {
	switch c {
	case BeetCandy:
		return "R"
	case MilkBread:
		return "Y"
	case Nucleons:
		return "G"
	case OilPetrol:
		return "B"
	default:
		return "*"
	}
}

// TimeLine represents a timeline that flows from one side of a PathCard to the other.
type TimeLine struct {
	Card *Card     // Card that this TimeLine is on
	Slot int       // 0..3
	From *TimeLine // Nil for Start cards, otherwise a link to a timeline on another card
	To   *TimeLine // Nil for Finish cards, otherwise a link to a timeline on another card
}

type Waypoint struct {
	Commodity Commodity
	// X and Y are the coordinates of the waypoint on the card.
	// The origin is the top left corner of the un-rotated card.
	// We need these coordinates to draw the token on the board.
	// More accurately, we need the coordinates of the center of the token so that we can erase it.
	X, Y float64
}
