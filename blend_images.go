package main

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

// calculateDimensions takes the width and height of two images and returns the best width and height to resize them to.
func calculateDimensions(width1, height1, width2, height2 uint) (newWidth, newHeight uint) {
	// Calculate the aspect ratios of both images
	aspectRatio1 := float64(width1) / float64(height1)
	aspectRatio2 := float64(width2) / float64(height2)

	// Determine which image is smaller in each dimension
	if width1 < width2 {
		newWidth = width1
	} else {
		newWidth = width2
	}

	if height1 < height2 {
		newHeight = height1
	} else {
		newHeight = height2
	}

	// Adjust the dimensions to maintain aspect ratio
	if aspectRatio1 < aspectRatio2 {
		// Adjust width for the first image
		newWidth = uint(float64(newHeight) * aspectRatio1)
	} else {
		// Adjust height for the second image
		newHeight = uint(float64(newWidth) / aspectRatio2)
	}

	return newWidth, newHeight
}

func resizeImages(image1 image.Image, image2 image.Image) (image.Image, image.Image) {
	image1Bounds := image1.Bounds()
	image2Bounds := image2.Bounds()
	width1, height1 := image1Bounds.Max.X, image1Bounds.Max.Y
	width2, height2 := image2Bounds.Max.X, image2Bounds.Max.Y
	newWidth, newHeight := calculateDimensions(uint(width1), uint(height1), uint(width2), uint(height2))
	return resize.Resize(newWidth, newHeight, image1, resize.Lanczos3), resize.Resize(newWidth, newHeight, image2, resize.Lanczos3)

}

func blendImages(image1 image.Image, image2 image.Image, modFunction func(pixel1 color.RGBA, pixel2 color.RGBA) color.RGBA) image.Image {
	img1, img2 := resizeImages(image1, image2)
	bounds := img1.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel1 := img1.At(x, y)
			pixel2 := img2.At(x, y)
			rgba := modFunction(color.RGBAModel.Convert(pixel1).(color.RGBA), color.RGBAModel.Convert(pixel2).(color.RGBA))
			newImage.Set(x, y, rgba)
		}
	}
	return newImage
}
func averagePixel(pixel1 color.RGBA, pixel2 color.RGBA) color.RGBA {
	red := (pixel1.R + pixel2.R) / 2
	green := (pixel1.G + pixel2.G) / 2
	blue := (pixel1.B + pixel2.B) / 2

	return color.RGBA{R: red, G: green, B: blue, A: 255}
}

func replaceHue(pixel1 color.RGBA, pixel2 color.RGBA) color.RGBA {
	hue1 := GetHueRatio(pixel1)
	light2 := GetTotalLight(pixel2)
	return ApplyLightToHue(hue1, light2)
}

func replaceBrightPixel(pixel1 color.RGBA, pixel2 color.RGBA) color.RGBA {
	if getTotalLight(pixel1) > 400 {
		return pixel2
	}
	return pixel1
}

func ReplaceBrightPixels(img1 image.Image, img2 image.Image) image.Image {
	return blendImages(img1, img2, replaceBrightPixel)
}
func BlendImages(img1 image.Image, img2 image.Image) image.Image {
	return blendImages(img1, img2, averagePixel)
}

func ReplaceHue(img1 image.Image, img2 image.Image) image.Image {
	return blendImages(img1, img2, replaceHue)
}
