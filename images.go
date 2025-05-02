package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// generateBirdImage creates a bird image
func generateBirdImage() *ebiten.Image {
	// Create a new image for the bird
	img := ebiten.NewImage(birdSize, birdSize)

	// Draw bird body (yellow circle with details)
	birdCenterX := float64(birdSize) / 2
	birdCenterY := float64(birdSize) / 2
	birdRadius := float64(birdSize) / 2

	// Draw main body (yellow circle)
	for y := 0; y < birdSize; y++ {
		for x := 0; x < birdSize; x++ {
			dx := float64(x) - birdCenterX
			dy := float64(y) - birdCenterY
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance < birdRadius {
				// Main body color (yellow)
				img.Set(x, y, birdColor)

				// Add shading for 3D effect
				if distance > birdRadius*0.8 {
					// Darker at the edges
					img.Set(x, y, color.RGBA{220, 220, 0, 255})
				}

				// Add belly (lighter yellow)
				if dy > 0 && dx < 0 {
					img.Set(x, y, color.RGBA{255, 255, 128, 255})
				}
			}
		}
	}

	// Draw eye (black circle with white highlight)
	eyeX := int(birdCenterX + birdRadius*0.3)
	eyeY := int(birdCenterY - birdRadius*0.2)
	eyeRadius := birdRadius * 0.2

	for y := 0; y < birdSize; y++ {
		for x := 0; x < birdSize; x++ {
			dx := float64(x) - float64(eyeX)
			dy := float64(y) - float64(eyeY)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance < eyeRadius {
				img.Set(x, y, color.Black)

				// Add white highlight
				if dx < 0 && dy < 0 && distance < eyeRadius*0.5 {
					img.Set(x, y, color.White)
				}
			}
		}
	}

	// Draw beak (orange triangle)
	beakX := int(birdCenterX + birdRadius*0.7)
	beakY := int(birdCenterY)
	beakWidth := int(birdRadius * 0.6)
	beakHeight := int(birdRadius * 0.4)

	for y := beakY - beakHeight/2; y <= beakY+beakHeight/2; y++ {
		if y < 0 || y >= birdSize {
			continue
		}

		// Calculate how far across the beak we should draw at this y coordinate
		yRatio := float64(y-(beakY-beakHeight/2)) / float64(beakHeight)
		xEnd := beakX + beakWidth

		// Triangle shape
		if yRatio <= 0.5 {
			// Top half of beak
			yRatioAdjusted := yRatio * 2
			xStart := beakX + int(float64(beakWidth)*(1-yRatioAdjusted))
			for x := xStart; x < xEnd && x < birdSize; x++ {
				img.Set(x, y, birdBeakColor)
			}
		} else {
			// Bottom half of beak
			yRatioAdjusted := (yRatio - 0.5) * 2
			xStart := beakX + int(float64(beakWidth)*yRatioAdjusted)
			for x := xStart; x < xEnd && x < birdSize; x++ {
				img.Set(x, y, birdBeakColor)
			}
		}
	}

	// Draw wing (yellow oval)
	wingCenterX := int(birdCenterX - birdRadius*0.4)
	wingCenterY := int(birdCenterY + birdRadius*0.1)
	wingRadiusX := birdRadius * 0.5
	wingRadiusY := birdRadius * 0.3

	for y := 0; y < birdSize; y++ {
		for x := 0; x < birdSize; x++ {
			dx := (float64(x) - float64(wingCenterX)) / wingRadiusX
			dy := (float64(y) - float64(wingCenterY)) / wingRadiusY
			distance := dx*dx + dy*dy

			if distance < 1 {
				// Wing color (slightly darker yellow)
				img.Set(x, y, color.RGBA{220, 220, 0, 255})

				// Add detail lines
				if int(float64(y-wingCenterY)*5)%4 == 0 {
					img.Set(x, y, color.RGBA{200, 200, 0, 255})
				}
			}
		}
	}

	return img
}

// generatePipeImage creates a pipe body image
func generatePipeImage() *ebiten.Image {
	// Create a new image for the pipe
	img := ebiten.NewImage(pipeWidth, 100) // Height will be scaled when drawing

	// Draw pipe body (green rectangle with texture)
	for y := 0; y < 100; y++ {
		for x := 0; x < pipeWidth; x++ {
			// Base pipe color
			img.Set(x, y, pipeColor)

			// Add vertical texture lines
			if x%15 == 0 || x%15 == 1 {
				img.Set(x, y, pipeBorderColor)
			}

			// Add horizontal texture lines
			if y%20 == 0 || y%20 == 1 {
				img.Set(x, y, pipeBorderColor)
			}

			// Add edge shading
			if x < 3 || x >= pipeWidth-3 {
				img.Set(x, y, pipeBorderColor)
			}
		}
	}

	return img
}

// generatePipeCapImage creates a pipe cap image
func generatePipeCapImage() *ebiten.Image {
	// Create a new image for the pipe cap
	capWidth := pipeWidth + 10
	capHeight := 20
	img := ebiten.NewImage(capWidth, capHeight)

	// Draw pipe cap (darker green rectangle with texture)
	for y := 0; y < capHeight; y++ {
		for x := 0; x < capWidth; x++ {
			// Base cap color
			img.Set(x, y, pipeBorderColor)

			// Add texture
			if (x+y)%5 == 0 {
				img.Set(x, y, color.RGBA{0, 80, 0, 255}) // Darker texture
			}

			// Add highlight on top
			if y < 3 {
				img.Set(x, y, color.RGBA{0, 150, 0, 255}) // Lighter green
			}

			// Add shadow at bottom
			if y > capHeight-4 {
				img.Set(x, y, color.RGBA{0, 60, 0, 255}) // Darker green
			}
		}
	}

	return img
}

// generateBackgroundImage creates a background image
func generateBackgroundImage() *ebiten.Image {
	// Create a new image for the background
	img := ebiten.NewImage(screenWidth, screenHeight)

	// Draw sky gradient
	for y := 0; y < screenHeight; y++ {
		// Calculate sky color based on height (gradient from light blue to darker blue)
		ratio := float64(y) / float64(screenHeight)
		r := uint8(135 - ratio*30)
		g := uint8(206 - ratio*50)
		b := uint8(235)
		skyColor := color.RGBA{r, g, b, 255}

		for x := 0; x < screenWidth; x++ {
			img.Set(x, y, skyColor)
		}
	}

	// Draw distant mountains
	mountainColor := color.RGBA{100, 120, 150, 255}
	mountainHeight := 100

	for i := 0; i < 5; i++ {
		// Draw a mountain
		peakX := screenWidth * (i + 1) / 6
		for x := peakX - 150; x < peakX+150; x++ {
			if x < 0 || x >= screenWidth {
				continue
			}

			// Calculate mountain height at this x position
			distFromPeak := math.Abs(float64(x - peakX))
			height := int(float64(mountainHeight) * (1 - distFromPeak/150))

			for y := screenHeight - height; y < screenHeight; y++ {
				if y < 0 || y >= screenHeight {
					continue
				}

				// Add some texture to the mountains
				noise := rand.Intn(20) - 10
				adjustedColor := color.RGBA{
					uint8(int(mountainColor.R) + noise),
					uint8(int(mountainColor.G) + noise),
					uint8(int(mountainColor.B) + noise),
					255,
				}

				img.Set(x, y, adjustedColor)
			}
		}
	}

	return img
}

// generateGroundImage creates a ground image
func generateGroundImage() *ebiten.Image {
	// Create a new image for the ground
	groundHeight := 30
	img := ebiten.NewImage(screenWidth, groundHeight)

	// Draw ground (brown with texture)
	for y := 0; y < groundHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			// Base ground color
			img.Set(x, y, groundColor)

			// Add texture
			if (x+y)%4 == 0 {
				img.Set(x, y, color.RGBA{120, 60, 15, 255}) // Darker spots
			}

			// Add grass on top
			if y < 5 && rand.Intn(10) < 3 {
				grassHeight := rand.Intn(5) + 1
				for gy := 0; gy < grassHeight && y-gy >= 0; gy++ {
					img.Set(x, y-gy, grassColor)
				}
			}
		}
	}

	return img
}
