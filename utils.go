package main

import (
	"fmt"
	"image" // For JPEG images
	"image/color"

	// For PNG images
	"os"
)

func LoadImageFromFile(filepath string) (image.Image, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil, err
	}

	// img is now an image.Image object
	fmt.Println("Image successfully loaded into memory:", img.Bounds())
	return img, nil
}

func CreateNewImage(width int, height int, getPixelColor func(x, y int) color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Set the color of each pixel
			img.Set(x, y, getPixelColor(x, y))
		}
	}

	return img
}
