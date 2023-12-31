package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

// Centroid represents the center of a cluster.
type Centroid struct {
	R, G, B, A float64
}

func QuantizeColors(pixels []color.RGBA) []color.RGBA {
	// Example: Load or create your Pixel data here.
	// Set the number of centroids (n).
	n := 10 // for example, reduce to 5 colors

	// Run k-means clustering.
	centroids := kMeans(pixels, n)

	// Example: Output centroids (dominant colors).
	return centroidsToRGBA(centroids)
}

func kMeans(pixels []color.RGBA, n int) []Centroid {
	// Initialize n centroids randomly.
	centroids := initializeCentroids(pixels, n)

	var assignments []int
	for {
		// Assign each pixel to the nearest centroid.
		newAssignments := assignPixelsToCentroids(pixels, centroids)

		// Check for convergence (no change in assignments).
		if equal(assignments, newAssignments) {
			break
		}
		assignments = newAssignments

		// Recalculate centroids.
		centroids = recalculateCentroids(pixels, assignments, n)
	}

	return centroids
}

func initializeCentroids(pixels []color.RGBA, n int) []Centroid {
	rand.Seed(time.Now().UnixNano())
	centroids := make([]Centroid, n)

	for i := range centroids {
		p := pixels[rand.Intn(len(pixels))]
		centroids[i] = Centroid{R: float64(p.R), G: float64(p.G), B: float64(p.B), A: float64(p.A)}
	}

	return centroids
}

func assignPixelsToCentroids(pixels []color.RGBA, centroids []Centroid) []int {
	assignments := make([]int, len(pixels))

	for i, p := range pixels {
		minDistance := math.MaxFloat64
		for j, c := range centroids {
			distance := euclideanDistance(p, c)
			if distance < minDistance {
				minDistance = distance
				assignments[i] = j
			}
		}
	}

	return assignments
}

func recalculateCentroids(pixels []color.RGBA, assignments []int, n int) []Centroid {
	sums := make([]Centroid, n)
	counts := make([]int, n)

	for i, p := range pixels {
		centroidIndex := assignments[i]
		sums[centroidIndex].R += float64(p.R)
		sums[centroidIndex].G += float64(p.G)
		sums[centroidIndex].B += float64(p.B)
		sums[centroidIndex].A += float64(p.A)
		counts[centroidIndex]++
	}

	for i := range sums {
		if counts[i] > 0 {
			sums[i].R /= float64(counts[i])
			sums[i].G /= float64(counts[i])
			sums[i].B /= float64(counts[i])
			sums[i].A /= float64(counts[i])
		}
	}

	return sums
}

func euclideanDistance(p color.RGBA, c Centroid) float64 {
	return math.Sqrt(math.Pow(float64(p.R)-c.R, 2) + math.Pow(float64(p.G)-c.G, 2) +
		math.Pow(float64(p.B)-c.B, 2) + math.Pow(float64(p.A)-c.A, 2))
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func centroidToRGBA(c Centroid) color.RGBA {
	return color.RGBA{
		R: uint8(c.R),
		G: uint8(c.G),
		B: uint8(c.B),
		A: uint8(c.A),
	}
}

// centroidsToRGBA converts a slice of Centroid to a slice of color.RGBA.
func centroidsToRGBA(centroids []Centroid) []color.RGBA {
	var rgbaSlice []color.RGBA
	for _, c := range centroids {
		rgbaColor := color.RGBA{
			R: uint8(c.R),
			G: uint8(c.G),
			B: uint8(c.B),
			A: uint8(c.A),
		}
		rgbaSlice = append(rgbaSlice, rgbaColor)
	}
	return rgbaSlice
}

func FindClosestQuantizedColor(pixel color.RGBA, quantizedColors []color.RGBA) color.RGBA {
	minDistance := math.MaxFloat64
	var closestColor color.RGBA

	for _, color := range quantizedColors {
		distance := colorDistance(pixel, color)
		if distance < minDistance {
			minDistance = distance
			closestColor = color
		}
	}

	return closestColor
}

func colorDistance(c1, c2 color.RGBA) float64 {
	return math.Sqrt(
		math.Pow(float64(c1.R)-float64(c2.R), 2) +
			math.Pow(float64(c1.G)-float64(c2.G), 2) +
			math.Pow(float64(c1.B)-float64(c2.B), 2) +
			math.Pow(float64(c1.A)-float64(c2.A), 2),
	)
}
