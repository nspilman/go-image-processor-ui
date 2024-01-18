package main

import (
	"image"
	"image/color"
	"testing"
)

func getTwoTestImages(dimension int) (image.Image, image.Image) {
	img1 := CreateNewImage(dimension, dimension, func(height, width int) color.RGBA {
		return color.RGBA{
			R: 200,
			G: 200,
			B: 200,
			A: 255,
		}
	})
	img2 := CreateNewImage(dimension, dimension, func(height, width int) color.RGBA {
		return color.RGBA{
			R: 100,
			G: 100,
			B: 100,
			A: 255,
		}
	})
	return img1, img2
}

func getTestDimension() int {
	return 300
}

func BenchmarkModifyImageConcurrent(b *testing.B) {
	img, img2 := getTwoTestImages(getTestDimension())

	b.ResetTimer() // Reset the timer to exclude the setup time

	for i := 0; i < b.N; i++ {
		_ = blendImagesConcurrently(img, img2, replaceHue)
	}
}

func BenchmarkModifyImage(b *testing.B) {
	img, img2 := getTwoTestImages(getTestDimension())

	b.ResetTimer() // Reset the timer to exclude the setup time

	for i := 0; i < b.N; i++ {
		_ = blendImages(img, img2, replaceHue)
	}
}
