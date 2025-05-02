package main

import (
	"image/color"
)

// Audio constants
const (
	sampleRate      = 44100
	audioBufferSize = 1024
)

const (
	// Game area dimensions
	screenWidth  = 640
	screenHeight = 480

	// Game physics
	gravity     = 0.3
	jumpForce   = -5.0
	pipeSpeed   = 4
	pipeGap     = 100
	pipeSpacing = 120
	pipeWidth   = 80

	// Bird dimensions
	birdSize = 30

	// Game states
	stateMenu = iota
	statePlay
)

// Colors
var (
	backgroundColor = color.RGBA{135, 206, 235, 255} // Sky blue
	birdColor       = color.RGBA{255, 255, 0, 255}   // Yellow
	birdEyeColor    = color.RGBA{0, 0, 0, 255}       // Black
	birdBeakColor   = color.RGBA{255, 165, 0, 255}   // Orange
	pipeColor       = color.RGBA{0, 128, 0, 255}     // Green
	pipeBorderColor = color.RGBA{0, 100, 0, 255}     // Dark green
	groundColor     = color.RGBA{139, 69, 19, 255}   // Brown
	grassColor      = color.RGBA{34, 139, 34, 255}   // Forest green
	cloudColor      = color.RGBA{255, 255, 255, 240} // White with transparency
	textColor       = color.RGBA{255, 255, 255, 255} // White
)
