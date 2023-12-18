package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func encodeImage(img image.Image) string {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return ""
	}
	base64Image := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return "data:image/png;base64," + base64Image
}

func iAmABanana() string {
	file, err := os.Open("BlueSailboats.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Step 3: Decode the image.
	img, err := jpeg.Decode(file)
	return encodeImage(img)
}
