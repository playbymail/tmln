// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

// tileAndCropImage is used to extract the images for cards and tokens from the supplied PDF file.
//
// You must fetch the PDF from the Crab Fragment Labs website (https://crabfragmentlabs.com/stage-blood).
func tileAndCropImage(img *image.RGBA, rows, cols int, topMargin, bottomMargin, leftMargin, rightMargin int) (*image.RGBA, [][]*image.RGBA) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	black := color.RGBA{0, 0, 0, 255}

	// draw the borders
	drawLine(img, 0+leftMargin, 0, 0+leftMargin, height, 1, black)              // left border
	drawLine(img, 0, height-bottomMargin, width, height-bottomMargin, 1, black) // bottom border
	drawLine(img, width-rightMargin, 0, width-rightMargin, height, 1, black)    // right border
	drawLine(img, 0, 0+topMargin, width, 0+topMargin, 1, black)                 // top border

	// Prepare a slice to store the cropped tiles
	var croppedTiles [][]*image.RGBA

	// Calculate the drawable area inside the margins
	drawableWidth := width - (leftMargin + rightMargin)
	drawableHeight := height - (topMargin + bottomMargin)
	//log.Printf("tile: width %6d, height %6d, drawableWidth %6d, drawableHeight %6d\n", width, height, drawableWidth, drawableHeight)

	// Calculate the width and height of each tile in the drawable area
	tileWidth := drawableWidth / cols
	tileHeight := drawableHeight / rows
	//log.Printf("tile: tileWidth %6d, tileHeight %6d\n", tileWidth, tileHeight)

	// Draw the horizontal lines between the tiles
	for i := 0; i <= rows; i++ {
		x1, y1 := 0, topMargin+i*tileHeight
		x2, y2 := width, topMargin+i*tileHeight
		drawLine(img, x1, y1, x2, y2, 3, black)
	}

	// Draw the vertical lines between the tiles
	for i := 0; i <= cols; i++ {
		x1, y1 := leftMargin+i*tileWidth, 0
		x2, y2 := leftMargin+i*tileWidth, height
		drawLine(img, x1, y1, x2, y2, 3, black)
	}

	// Now, crop each tile and store it in the slice
	for row := 0; row < rows; row++ {
		croppedTiles = append(croppedTiles, []*image.RGBA{})
		for col := 0; col < cols; col++ {
			// Calculate the bounds of each tile
			minX := leftMargin + col*tileWidth
			minY := topMargin + row*tileHeight
			maxX := minX + tileWidth
			maxY := minY + tileHeight

			// Crop the tile using the bounds
			tile := img.SubImage(image.Rect(minX, minY, maxX, maxY)).(*image.RGBA)

			// Append the cropped tile to the slice
			croppedTiles[row] = append(croppedTiles[row], tile)
		}
	}

	// Return the updated image and the slice of cropped tiles
	return img, croppedTiles
}

// helper functions

// drawLine draws a line between two points (x1, y1) and (x2, y2) with a given width and color, clipped to the image bounds.
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, width int, c color.RGBA) {
	if width <= 0 {
		width = 1 // Ensure width is at least 1 pixel
	}

	// Clip the start and end points to the image bounds
	imgBounds := img.Bounds()
	x1 = clip(x1, imgBounds.Min.X, imgBounds.Max.X-1)
	y1 = clip(y1, imgBounds.Min.Y, imgBounds.Max.Y-1)
	x2 = clip(x2, imgBounds.Min.X, imgBounds.Max.X-1)
	y2 = clip(y2, imgBounds.Min.Y, imgBounds.Max.Y-1)

	// Use Bresenham's line algorithm to draw the line
	dx, dy := abs(x2-x1), abs(y2-y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	for {
		// Draw a "thicker" line by drawing additional pixels around the main line
		for i := -width / 2; i <= width/2; i++ {
			for j := -width / 2; j <= width/2; j++ {
				if x1+i >= imgBounds.Min.X && x1+i < imgBounds.Max.X && y1+j >= imgBounds.Min.Y && y1+j < imgBounds.Max.Y {
					img.Set(x1+i, y1+j, c)
				}
			}
		}

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// abs returns the absolute value of an integer.
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// clip ensures that a value stays within the given min and max range.
func clip(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// cropImage crops an image to the specified rectangle
func cropImage(img image.Image, left, upper, right, lower int) image.Image {
	cropped := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(left, upper, right, lower))
	return cropped
}

// saveImageAsPNG saves the image as a PNG file
func saveImageAsPNG(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Printf("error: saving image: %v\n", err)
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
