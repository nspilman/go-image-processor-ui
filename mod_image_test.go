package main

import (
	"fmt"
	"image/color"
	"testing"
)

func TestGetTotalLight_returnsExpected_whenNoAlphaAndFullAlpha(t *testing.T) {
	no_alpha_result := getTotalLight(color.RGBA{100, 100, 100, 0})
	full_alpha_result := getTotalLight(color.RGBA{100, 100, 100, 255})

	expected := uint16(300)
	if no_alpha_result != expected || full_alpha_result != expected {
		t.Errorf("Expected %d, got %d and %d", expected, full_alpha_result, no_alpha_result)
	}
}

func TestMaxGuardAbove(t *testing.T) {
	result := maxGuard(256)
	expected := uint8(255)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestMaxGuardBelow(t *testing.T) {
	result := maxGuard(254)
	expected := uint8(254)
	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}

func TestGetHueRatio_ReturnsFullRedRatio_WhenFullRedPixelPassed(t *testing.T) {
	RedPixel := color.RGBA{255, 0, 0, 255}
	result := getHueRatio(RedPixel)
	fmt.Println(result)
	if result.B != 0 || result.R != 1 || result.G != 0 {
		t.Errorf("error")
	}
}
