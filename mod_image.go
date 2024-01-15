package main

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/nfnt/resize"
)

const MAX = 255

func modifyImage(img image.Image, modFunction func(pixel color.RGBA) color.RGBA) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			rgba := modFunction(color.RGBAModel.Convert(pixel).(color.RGBA))
			newImage.Set(x, y, rgba)
		}
	}
	return newImage
}

func modifyImageConcurrently(img image.Image, modFunction func(pixel color.RGBA) color.RGBA) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	var mu sync.Mutex
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				pixel := img.At(x, y)
				rgba := modFunction(color.RGBAModel.Convert(pixel).(color.RGBA))
				mu.Lock()
				newImage.Set(x, y, rgba)
				mu.Unlock()
			}(x, y)
		}
	}
	wg.Wait()
	return newImage
}

func getTotalLight(pixel color.RGBA) uint16 {
	return uint16(pixel.R) + uint16(pixel.G) + uint16(pixel.B)
}

func GetTotalLight(pixel color.RGBA) uint16 {
	return getTotalLight(pixel)
}

type PixelRatio struct {
	R float32
	G float32
	B float32
}

type Pixel struct {
	R, G, B, A uint8
}

func getHueRatio(pixel color.RGBA) PixelRatio {
	light := getTotalLight(pixel)
	if light == 0 {
		return PixelRatio{
			R: float32(1 / 3),
			G: float32(1 / 3),
			B: float32(1 / 3),
		}
	}
	ratio := PixelRatio{
		R: float32(pixel.R) / float32(light),
		G: float32(pixel.G) / float32(light),
		B: float32(pixel.B) / float32(light),
	}
	return ratio
}

func GetHueRatio(pixel color.RGBA) PixelRatio {
	return getHueRatio(pixel)
}

func invertLight(pixel color.RGBA) color.RGBA {
	light := 755 - getTotalLight(pixel)
	hue := getHueRatio(pixel)
	return applyLightToHue(hue, light)
}

func InvertLight(img image.Image) image.Image {
	return modifyImage(img, invertLight)
}

func normalizeLight(pixel color.RGBA) color.RGBA {
	hue := getHueRatio(pixel)
	light := uint16(500)
	return applyLightToHue(hue, light)
}

func NormalizeLight(img image.Image) image.Image {
	return modifyImage(img, normalizeLight)
}

func maxGuard(color float32) uint8 {
	if color > MAX {
		return MAX
	}
	return uint8(color)
}

func applyLightToHue(hue PixelRatio, light uint16) color.RGBA {
	red := maxGuard(hue.R * float32(light))
	green := maxGuard(hue.G * float32(light))
	blue := maxGuard(hue.B * float32(light))
	return color.RGBA{
		R: red,
		G: green,
		B: blue,
		A: MAX,
	}
}

func ApplyLightToHue(hue PixelRatio, light uint16) color.RGBA {
	return applyLightToHue(hue, light)
}

func FlattenImage(img image.Image) []color.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var flattened []color.RGBA
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			flattened = append(flattened, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	return flattened
}

func QuantizeImage(imgs []image.Image) image.Image {
	img := imgs[0]
	bounds := img.Bounds()
	width1, height1 := bounds.Max.X, bounds.Max.Y

	resized := resize.Resize(uint(width1/100), uint(height1/100), img, resize.Lanczos3)

	flatted := FlattenImage(resized)

	colors := QuantizeColors(flatted)

	getClosestPixelColor := func(pixel color.RGBA) color.RGBA {
		return FindClosestQuantizedColor(pixel, colors)
	}
	return modifyImage(img, getClosestPixelColor)
}

func roundToNearest(number float64, interval int) int {
	divisor := float64(interval)
	rounded := math.Floor(number/divisor) * divisor
	remainder := number - rounded

	if remainder >= divisor/2 {
		rounded += divisor
	}

	return int(rounded)
}

func quantizeLight(pixel color.RGBA) color.RGBA {
	var newLight int
	quantizedLight := roundToNearest(float64(getTotalLight(pixel)), 100)
	if quantizedLight > (255 * 3) {
		newLight = 255 * 3
	} else {
		newLight = quantizedLight
	}
	return applyLightToHue(getHueRatio((pixel)), uint16(newLight))
}

func QuantizeLight(imgs []image.Image) image.Image {
	img := imgs[0]
	return modifyImage(img, quantizeLight)
}
