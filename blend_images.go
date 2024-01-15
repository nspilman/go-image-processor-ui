package main

import (
	"image"
	"image/color"
	"os"
	"sync"

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

func blendImagesConcurrently(image1, image2 image.Image, modFunction func(pixel1, pixel2 color.RGBA) color.RGBA) image.Image {
	img1, img2 := resizeImages(image1, image2) // Assuming resizeImages is defined elsewhere
	bounds := img1.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup

	// Determine the number of goroutines to use
	numGoroutines := 16 // For example, can be tuned based on the environment
	rowsPerGoroutine := height / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startRow, endRow int) {
			defer wg.Done()
			for y := startRow; y < endRow; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					pixel1 := img1.At(x, y)
					pixel2 := img2.At(x, y)
					rgba := modFunction(color.RGBAModel.Convert(pixel1).(color.RGBA), color.RGBAModel.Convert(pixel2).(color.RGBA))
					newImage.Set(x, y, rgba)
				}
			}
		}(i*rowsPerGoroutine, min((i+1)*rowsPerGoroutine, height))
	}
	wg.Wait()
	return newImage
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func averagePixel(pixel1 color.RGBA, pixel2 color.RGBA) color.RGBA {
	hue1, hue2 := getHueRatio(pixel1), getHueRatio(pixel2)
	newHue := PixelRatio{R: (hue1.R + hue2.R) / 2, G: (hue1.G + hue2.G) / 2, B: (hue1.B + hue2.B) / 2}
	newLight := getTotalLight(pixel1) + getTotalLight(pixel2)
	return applyLightToHue(newHue, newLight)
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

func ReplaceBrightPixels(imgs []image.Image) image.Image {
	img1, img2 := expectTwoImages(imgs)
	return blendImages(img1, img2, replaceBrightPixel)
}

func expectTwoImages(imgs []image.Image) (image.Image, image.Image) {
	if len(imgs) != 2 {
		os.Exit(1)
	}
	return imgs[0], imgs[1]
}
func BlendImages(imgs []image.Image) image.Image {
	img1, img2 := expectTwoImages(imgs)
	return blendImages(img1, img2, averagePixel)
}

func ReplaceHue(imgs []image.Image) image.Image {
	img1, img2 := expectTwoImages(imgs)
	return blendImagesConcurrently(img1, img2, replaceHue)
}
