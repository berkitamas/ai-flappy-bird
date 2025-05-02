package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// NewGame creates a new game instance
func NewGame() *Game {
	// Create initial clouds
	initialClouds := make([]Cloud, 5)
	for i := range initialClouds {
		initialClouds[i] = Cloud{
			x:      float64(rand.Intn(screenWidth)),
			y:      float64(rand.Intn(screenHeight / 3)),
			width:  float64(rand.Intn(80) + 40),
			height: float64(rand.Intn(30) + 20),
			speed:  float64(rand.Intn(2) + 1),
		}
	}

	// Generate images
	birdImg := generateBirdImage()
	pipeImg := generatePipeImage()
	pipeCapImg := generatePipeCapImage()
	backgroundImg := generateBackgroundImage()
	groundImg := generateGroundImage()

	// Initialize audio context
	audioContext := audio.NewContext(sampleRate)

	// Generate audio players
	menuMusic := generateMenuMusic(audioContext)
	jumpSound := generateJumpSound(audioContext)
	scoreSound := generateScoreSound(audioContext)
	collisionSound := generateCollisionSound(audioContext)

	// Start playing menu music
	menuMusic.Play()

	return &Game{
		bird: Bird{
			x:         float64(screenWidth) / 4,
			y:         float64(screenHeight) / 2,
			width:     birdSize,
			height:    birdSize,
			wingAngle: 0,
		},
		pipes:        []Pipe{},
		clouds:       initialClouds,
		groundOffset: 0,
		score:        0,
		isGameOver:   false,
		tick:         0,
		state:        stateMenu,
		menuOption:   0,
		demoMode:     true,
		keyStates:    make(map[ebiten.Key]bool),

		// Image resources
		birdImg:       birdImg,
		pipeImg:       pipeImg,
		pipeCapImg:    pipeCapImg,
		backgroundImg: backgroundImg,
		groundImg:     groundImg,

		// Audio resources
		audioContext:   audioContext,
		menuMusic:      menuMusic,
		jumpSound:      jumpSound,
		scoreSound:     scoreSound,
		collisionSound: collisionSound,
	}
}

// isKeyJustPressed checks if a key was just pressed this frame
func (g *Game) isKeyJustPressed(key ebiten.Key) bool {
	wasPressed := g.keyStates[key]
	isPressed := ebiten.IsKeyPressed(key)
	g.keyStates[key] = isPressed
	return isPressed && !wasPressed
}

// Update updates the game state
func (g *Game) Update() error {
	g.tick++

	// Update common elements only if the game is not over
	if !g.isGameOver {
		g.updateBackground()
	}

	// Handle state-specific updates
	switch g.state {
	case stateMenu:
		g.updateMenu()
	case statePlay:
		g.updateGame()
	}

	return nil
}

// updateBackground updates background elements like clouds and ground
func (g *Game) updateBackground() {
	// Update ground offset for scrolling effect
	g.groundOffset -= pipeSpeed
	if g.groundOffset <= -50 { // Reset when it scrolls one tile width
		g.groundOffset = 0
	}

	// Update clouds
	for i := range g.clouds {
		g.clouds[i].x -= g.clouds[i].speed

		// If cloud moves off-screen, reset it to the right side
		if g.clouds[i].x+g.clouds[i].width < 0 {
			g.clouds[i].x = float64(screenWidth)
			g.clouds[i].y = float64(rand.Intn(screenHeight / 3))
			g.clouds[i].width = float64(rand.Intn(80) + 40)
			g.clouds[i].height = float64(rand.Intn(30) + 20)
			g.clouds[i].speed = float64(rand.Intn(2) + 1)
		}
	}
}

// updateMenu handles menu state updates
func (g *Game) updateMenu() {
	// Handle menu input
	if g.isKeyJustPressed(ebiten.KeyDown) {
		g.menuOption = (g.menuOption + 1) % 2 // Cycle through menu options
	}
	if g.isKeyJustPressed(ebiten.KeyUp) {
		g.menuOption = (g.menuOption - 1 + 2) % 2 // Cycle through menu options
	}
	if g.isKeyJustPressed(ebiten.KeySpace) || g.isKeyJustPressed(ebiten.KeyEnter) {
		if g.menuOption == 0 { // Start Game option
			g.state = statePlay
			g.Reset()
			g.demoMode = false

			// Stop menu music when transitioning to play mode
			g.menuMusic.Pause()
		} else { // Exit option
			// In a real application, you would exit here
			// Since we can't call os.Exit in this environment, we'll just reset
			g.Reset()
		}
	}

	// Update demo game in background
	g.updateGameplay(true)
}

// updateGame handles gameplay state updates
func (g *Game) updateGame() {
	// Handle game input
	if g.isKeyJustPressed(ebiten.KeyEscape) {
		g.state = stateMenu
		g.Reset()
		g.demoMode = true

		// Start menu music when transitioning back to menu mode
		g.menuMusic.Rewind()
		g.menuMusic.Play()
	}

	if g.isGameOver {
		// Handle restart input
		if g.isKeyJustPressed(ebiten.KeyR) {
			g.Reset()
		}
		return
	}

	// Handle jump input
	if g.isKeyJustPressed(ebiten.KeySpace) {
		g.bird.velocity = jumpForce
		g.jumpSound.Rewind()
		g.jumpSound.Play()
	}

	// Update gameplay elements
	g.updateGameplay(false)
}

// Reset resets the game state
func (g *Game) Reset() {
	// Save current state and demoMode
	currentState := g.state
	demoMode := g.demoMode

	// Reset bird position and velocity
	g.bird.x = float64(screenWidth) / 4
	g.bird.y = float64(screenHeight) / 2
	g.bird.velocity = 0

	// Clear pipes
	g.pipes = []Pipe{}

	// Reset score and game over state
	g.score = 0
	g.isGameOver = false

	// Reset tick counter
	g.tick = 0

	// Restore state and demoMode
	g.state = currentState
	g.demoMode = demoMode

	// Clear key states
	g.keyStates = make(map[ebiten.Key]bool)
}

// Layout implements ebiten.Game's Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
