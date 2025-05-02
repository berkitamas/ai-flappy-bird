package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Bird represents the player character
type Bird struct {
	x, y      float64
	velocity  float64
	width     int
	height    int
	wingAngle float64 // For wing flapping animation
}

// Cloud represents a background cloud
type Cloud struct {
	x, y   float64
	width  float64
	height float64
	speed  float64
}

// Pipe represents an obstacle
type Pipe struct {
	x        int
	gapStart int
}

// Game represents the game state
type Game struct {
	bird         Bird
	pipes        []Pipe
	clouds       []Cloud
	groundOffset float64
	score        int
	isGameOver   bool
	tick         int
	state        int
	menuOption   int
	demoMode     bool
	keyStates    map[ebiten.Key]bool

	// Image resources
	birdImg       *ebiten.Image
	pipeImg       *ebiten.Image
	pipeCapImg    *ebiten.Image
	backgroundImg *ebiten.Image
	groundImg     *ebiten.Image

	// Audio resources
	audioContext   *audio.Context
	menuMusic      *audio.Player
	jumpSound      *audio.Player
	scoreSound     *audio.Player
	collisionSound *audio.Player
}
