package main

import (
	"image"
	"image/color"
)

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

func getTotalLight(pixel color.RGBA) uint16 {
	return uint16(pixel.R) + uint16(pixel.G) + uint16(pixel.B)
}

type PixelRatio struct {
	R float32
	G float32
	B float32
}

func getHueRatio(pixel color.RGBA) PixelRatio {
	light := getTotalLight(pixel)
	ratio := PixelRatio{
		R: float32(pixel.R) / float32(light),
		G: float32(pixel.G) / float32(light),
		B: float32(pixel.B) / float32(light),
	}
	return ratio
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
	light := uint16(750)
	return applyLightToHue(hue, light)
}

func NormalizeLight(img image.Image) image.Image {
	return modifyImage(img, normalizeLight)
}

func maxGuard(color float32) uint8 {
	if color > 255 {
		return 255
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
		A: 255,
	}
}
