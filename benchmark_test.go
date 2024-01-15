package main

import (
	"image/color"
	"testing"
)

func BenchmarkModifyImageConcurrent(b *testing.B) {
	dimension := 250

	// Create a sample image

	img := CreateNewImage(dimension, dimension, func(height, width int) color.RGBA {
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

	b.ResetTimer() // Reset the timer to exclude the setup time

	for i := 0; i < b.N; i++ {
		_ = blendImagesConcurrently(img, img2, replaceHue)
	}
}

func BenchmarkModifyImage(b *testing.B) {
	// Create a sample image
	dimension := 250

	img := CreateNewImage(dimension, dimension, func(height, width int) color.RGBA {
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

	b.ResetTimer() // Reset the timer to exclude the setup time

	for i := 0; i < b.N; i++ {
		_ = blendImages(img, img2, replaceHue)
	}
}
