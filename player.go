// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package tmln

type Player struct {
	Handle   string
	Tokens   map[Commodity]int
	Gold     int
	Location *TimeLine
}
