// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"log"
	"path/filepath"
)

// extractTimeLineCardsFromPDF extracts the images for cards from a PDF file and saves them as individual PNG files
func extractTimeLineCardsFromPDF(doc *fitz.Document, outputDir string) error {
	const cardWidth, leftMargin, rightMargin = 750, 149, 149
	const cardHeight, topMargin, bottomMargin = 1048, 76, 76
	const rows, cols = 3, 3
	log.Printf("card: height %6d (%3d): width %6d (%3d)", cardHeight, topMargin, cardWidth, leftMargin)

	// Iterate over each page
	for pageNum := 0; pageNum < doc.NumPage(); pageNum++ {
		img, err := doc.Image(pageNum)
		if err != nil {
			log.Printf("error: getting image from page: %v\n", err)
			return err
		}

		// Get image dimensions
		width := img.Bounds().Max.X
		height := img.Bounds().Max.Y
		log.Printf("page: %2d: %6dx%6d", pageNum+1, width, height)

		img, tiles := tileAndCropImage(img, rows, cols, topMargin, bottomMargin, leftMargin, rightMargin)
		tileImagePath := filepath.Join(outputDir, fmt.Sprintf("page_%02d.png", pageNum+1))
		if err := saveImageAsPNG(img, tileImagePath); err != nil {
			log.Printf("error saving %s: %v", tileImagePath, err)
			return err
		}
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				tile := tiles[row][col]
				tileImagePath := filepath.Join(outputDir, fmt.Sprintf("page_%02d_row_%02d_col_%02d.png", pageNum+1, row+1, col+1))
				if err := saveImageAsPNG(tile, tileImagePath); err != nil {
					log.Printf("error saving %s: %v", tileImagePath, err)
					return err
				}
			}
		}
	}

	log.Printf("pdf: timeline cards have been successfully extracted and saved.\n")

	return nil
}
