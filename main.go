package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

// generateSineWave generates a sine wave audio sample
func generateSineWave(freq float64, duration time.Duration) []byte {
	// Calculate the number of samples
	numSamples := int(sampleRate * duration.Seconds())

	// Create a buffer for the samples
	samples := make([]byte, numSamples*2) // 16-bit samples = 2 bytes per sample

	// Generate the sine wave
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		amplitude := math.Sin(2 * math.Pi * freq * t)

		// Convert to 16-bit PCM
		sample := int16(amplitude * 32767)

		// Write the sample to the buffer (little endian)
		samples[i*2] = byte(sample)
		samples[i*2+1] = byte(sample >> 8)
	}

	return samples
}

// addWAVHeader adds a WAV header to raw PCM data
func addWAVHeader(pcmData []byte) []byte {
	// WAV header constants
	const (
		headerSize    = 44
		formatPCM     = 1
		numChannels   = 1
		bitsPerSample = 16
	)

	// Calculate sizes
	dataSize := len(pcmData)
	fileSize := headerSize + dataSize - 8 // File size minus the first 8 bytes

	// Create header buffer
	header := make([]byte, headerSize)

	// RIFF header
	copy(header[0:4], []byte("RIFF"))
	// File size minus the first 8 bytes (little endian)
	header[4] = byte(fileSize & 0xFF)
	header[5] = byte((fileSize >> 8) & 0xFF)
	header[6] = byte((fileSize >> 16) & 0xFF)
	header[7] = byte((fileSize >> 24) & 0xFF)
	// WAVE identifier
	copy(header[8:12], []byte("WAVE"))

	// fmt chunk
	copy(header[12:16], []byte("fmt "))
	// fmt chunk size (16 for PCM) (little endian)
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	// Audio format (1 for PCM) (little endian)
	header[20] = byte(formatPCM)
	header[21] = 0
	// Number of channels (little endian)
	header[22] = byte(numChannels)
	header[23] = 0
	// Sample rate (little endian)
	sampleRateInt := int(sampleRate)
	header[24] = byte(sampleRateInt & 0xFF)
	header[25] = byte((sampleRateInt >> 8) & 0xFF)
	header[26] = byte((sampleRateInt >> 16) & 0xFF)
	header[27] = byte((sampleRateInt >> 24) & 0xFF)
	// Byte rate = SampleRate * NumChannels * BitsPerSample/8 (little endian)
	byteRate := int(sampleRate) * numChannels * bitsPerSample / 8
	header[28] = byte(byteRate & 0xFF)
	header[29] = byte((byteRate >> 8) & 0xFF)
	header[30] = byte((byteRate >> 16) & 0xFF)
	header[31] = byte((byteRate >> 24) & 0xFF)
	// Block align = NumChannels * BitsPerSample/8 (little endian)
	blockAlign := numChannels * bitsPerSample / 8
	header[32] = byte(blockAlign)
	header[33] = byte(blockAlign >> 8)
	// Bits per sample (little endian)
	header[34] = byte(bitsPerSample)
	header[35] = 0

	// data chunk
	copy(header[36:40], []byte("data"))
	// Data size (little endian)
	header[40] = byte(dataSize & 0xFF)
	header[41] = byte((dataSize >> 8) & 0xFF)
	header[42] = byte((dataSize >> 16) & 0xFF)
	header[43] = byte((dataSize >> 24) & 0xFF)

	// Combine header and PCM data
	wavData := make([]byte, headerSize+dataSize)
	copy(wavData[0:headerSize], header)
	copy(wavData[headerSize:], pcmData)

	return wavData
}

// generateJumpSound creates a jump sound
func generateJumpSound(context *audio.Context) *audio.Player {
	// Generate a short upward frequency sweep
	duration := 200 * time.Millisecond
	numSamples := int(sampleRate * duration.Seconds())
	samples := make([]byte, numSamples*2) // 16-bit samples = 2 bytes per sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		progress := t / duration.Seconds()

		// Frequency sweep from 300Hz to 600Hz
		freq := 300 + 300*progress

		// Amplitude envelope (fade in and out)
		envelope := 1.0
		if progress < 0.1 {
			envelope = progress / 0.1 // Fade in
		} else if progress > 0.7 {
			envelope = (1.0 - progress) / 0.3 // Fade out
		}

		amplitude := math.Sin(2*math.Pi*freq*t) * envelope

		// Convert to 16-bit PCM
		sample := int16(amplitude * 32767)

		// Write the sample to the buffer (little endian)
		samples[i*2] = byte(sample)
		samples[i*2+1] = byte(sample >> 8)
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Create a WAV reader from the samples with header
	wavData := bytes.NewReader(wavBytes)
	stream, err := wav.Decode(context, wavData)
	if err != nil {
		log.Fatal("Failed to decode jump sound: ", err)
	}

	// Create a player for the sound
	player, err := context.NewPlayer(stream)
	if err != nil {
		log.Fatal("Failed to create jump sound player: ", err)
	}

	return player
}

// generateScoreSound creates a score sound
func generateScoreSound(context *audio.Context) *audio.Player {
	// Generate a short "ding" sound
	duration := 300 * time.Millisecond
	numSamples := int(sampleRate * duration.Seconds())
	samples := make([]byte, numSamples*2) // 16-bit samples = 2 bytes per sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		progress := t / duration.Seconds()

		// Two frequencies for a pleasant sound
		freq1 := 800.0
		freq2 := 1200.0

		// Amplitude envelope (fade in and out)
		envelope := 1.0
		if progress < 0.05 {
			envelope = progress / 0.05 // Fade in
		} else if progress > 0.7 {
			envelope = (1.0 - progress) / 0.3 // Fade out
		}

		amplitude := (math.Sin(2*math.Pi*freq1*t)*0.7 +
			math.Sin(2*math.Pi*freq2*t)*0.3) * envelope

		// Convert to 16-bit PCM
		sample := int16(amplitude * 32767)

		// Write the sample to the buffer (little endian)
		samples[i*2] = byte(sample)
		samples[i*2+1] = byte(sample >> 8)
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Create a WAV reader from the samples with header
	wavData := bytes.NewReader(wavBytes)
	stream, err := wav.Decode(context, wavData)
	if err != nil {
		log.Fatal("Failed to decode score sound: ", err)
	}

	// Create a player for the sound
	player, err := context.NewPlayer(stream)
	if err != nil {
		log.Fatal("Failed to create score sound player: ", err)
	}

	return player
}

// generateCollisionSound creates a collision sound
func generateCollisionSound(context *audio.Context) *audio.Player {
	// Generate a short "crash" sound
	duration := 500 * time.Millisecond
	numSamples := int(sampleRate * duration.Seconds())
	samples := make([]byte, numSamples*2) // 16-bit samples = 2 bytes per sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		progress := t / duration.Seconds()

		// Noise-based sound with decreasing frequency
		freq := 400 - 300*progress
		noise := rand.Float64()*2 - 1 // Random value between -1 and 1

		// Amplitude envelope (quick fade in, longer fade out)
		envelope := 1.0
		if progress < 0.02 {
			envelope = progress / 0.02 // Fade in
		} else if progress > 0.3 {
			envelope = (1.0 - progress) / 0.7 // Fade out
		}

		// Mix sine wave and noise
		amplitude := (math.Sin(2*math.Pi*freq*t)*0.3 + noise*0.7) * envelope

		// Convert to 16-bit PCM
		sample := int16(amplitude * 32767)

		// Write the sample to the buffer (little endian)
		samples[i*2] = byte(sample)
		samples[i*2+1] = byte(sample >> 8)
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Create a WAV reader from the samples with header
	wavData := bytes.NewReader(wavBytes)
	stream, err := wav.Decode(context, wavData)
	if err != nil {
		log.Fatal("Failed to decode collision sound: ", err)
	}

	// Create a player for the sound
	player, err := context.NewPlayer(stream)
	if err != nil {
		log.Fatal("Failed to create collision sound player: ", err)
	}

	return player
}

// generateMenuMusic creates background music for the menu
func generateMenuMusic(context *audio.Context) *audio.Player {
	// Generate a simple looping melody
	duration := 10 * time.Second // 10 second loop
	numSamples := int(sampleRate * duration.Seconds())
	samples := make([]byte, numSamples*2) // 16-bit samples = 2 bytes per sample

	// Define a simple melody (frequencies and durations)
	type Note struct {
		freq     float64
		duration float64 // in seconds
	}

	melody := []Note{
		{392.0, 0.5}, // G4
		{440.0, 0.5}, // A4
		{392.0, 0.5}, // G4
		{440.0, 0.5}, // A4
		{493.9, 0.5}, // B4
		{440.0, 0.5}, // A4
		{392.0, 1.0}, // G4
		{349.2, 0.5}, // F4
		{392.0, 0.5}, // G4
		{349.2, 0.5}, // F4
		{392.0, 0.5}, // G4
		{440.0, 0.5}, // A4
		{392.0, 0.5}, // G4
		{349.2, 1.0}, // F4
		{329.6, 0.5}, // E4
		{349.2, 0.5}, // F4
		{329.6, 0.5}, // E4
		{349.2, 0.5}, // F4
		{392.0, 0.5}, // G4
		{349.2, 0.5}, // F4
		{329.6, 1.0}, // E4
	}

	// Generate the melody
	currentTime := 0.0
	for _, note := range melody {
		startSample := int(currentTime * sampleRate)
		endSample := int((currentTime + note.duration) * sampleRate)

		if endSample > numSamples {
			endSample = numSamples
		}

		for i := startSample; i < endSample; i++ {
			t := float64(i) / sampleRate

			// Calculate progress within this note for envelope
			noteProgress := (t - currentTime) / note.duration

			// Simple ADSR envelope
			envelope := 1.0
			if noteProgress < 0.1 {
				envelope = noteProgress / 0.1 // Attack
			} else if noteProgress > 0.8 {
				envelope = (1.0 - noteProgress) / 0.2 // Release
			}

			// Generate the waveform (sine wave)
			amplitude := math.Sin(2*math.Pi*note.freq*t) * envelope * 0.5

			// Add a simple accompaniment (chord)
			amplitude += math.Sin(2*math.Pi*(note.freq/2)*t) * envelope * 0.2
			amplitude += math.Sin(2*math.Pi*(note.freq*1.5)*t) * envelope * 0.1

			// Convert to 16-bit PCM
			sample := int16(amplitude * 32767)

			// Write the sample to the buffer (little endian)
			samples[i*2] = byte(sample)
			samples[i*2+1] = byte(sample >> 8)
		}

		currentTime += note.duration
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Create a WAV reader from the samples with header
	wavData := bytes.NewReader(wavBytes)
	stream, err := wav.Decode(context, wavData)
	if err != nil {
		log.Fatal("Failed to decode menu music: ", err)
	}

	// Create a player for the music
	player, err := context.NewPlayer(stream)
	if err != nil {
		log.Fatal("Failed to create menu music player: ", err)
	}

	return player
}

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
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.isGameOver {
		g.bird.velocity = jumpForce
		// Play jump sound
		g.jumpSound.Rewind()
		g.jumpSound.Play()
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && g.isGameOver {
		g.Reset()
	}

	// Update actual gameplay
	g.updateGameplay(false)
}

// updateGameplay handles the actual game mechanics
func (g *Game) updateGameplay(isDemo bool) {
	if !g.isGameOver {
		// Update bird position
		g.bird.velocity += gravity
		g.bird.y += g.bird.velocity

		// Update bird wing animation
		g.bird.wingAngle = 0.3 * g.bird.velocity
		if g.bird.wingAngle > 0.5 {
			g.bird.wingAngle = 0.5
		} else if g.bird.wingAngle < -0.5 {
			g.bird.wingAngle = -0.5
		}

		// Generate new pipes
		if g.tick%pipeSpacing == 0 {
			gapStart := rand.Intn(screenHeight-pipeGap-100) + 50
			g.pipes = append(g.pipes, Pipe{
				x:        screenWidth,
				gapStart: gapStart,
			})
		}

		// Update pipes and check for score
		for i := 0; i < len(g.pipes); i++ {
			g.pipes[i].x -= pipeSpeed

			// Check if bird passed the pipe
			// Use a range check instead of exact equality to avoid missing the score point
			if g.pipes[i].x+pipeWidth <= int(g.bird.x) && g.pipes[i].x+pipeWidth > int(g.bird.x)-pipeSpeed {
				g.score++
				// Play score sound
				g.scoreSound.Rewind()
				g.scoreSound.Play()
			}

			// Remove pipes that are off-screen
			if g.pipes[i].x+pipeWidth < 0 {
				g.pipes = append(g.pipes[:i], g.pipes[i+1:]...)
				i--
			}
		}

		// AI control for demo mode
		if isDemo {
			g.controlBirdAI()
		}

		// Check for collisions
		if !isDemo {
			g.checkCollisions()
		} else {
			// In demo mode, we still want to detect collisions but handle them differently
			// Set isGameOver to true to stop the game mechanics, including background movement
			if g.wouldCollide() {
				g.isGameOver = true
				// Play collision sound
				g.collisionSound.Rewind()
				g.collisionSound.Play()
				// Reset the game after a short delay in demo mode
				if g.tick%60 == 0 { // Reset approximately every second
					g.Reset()
				}
			}
		}
	}
}

// wouldCollide checks if the bird would collide with pipes or boundaries without setting game over
func (g *Game) wouldCollide() bool {
	// Calculate bird center and radius for circle-based collision
	birdCenterX := g.bird.x + float64(g.bird.width)/2
	birdCenterY := g.bird.y + float64(g.bird.height)/2
	birdRadius := float64(g.bird.width) / 2

	// Check if bird hit the top boundary
	if birdCenterY-birdRadius <= 0 {
		return true
	}

	// Check if bird hit the ground
	groundY := float64(screenHeight) - 30.0 // 30 is the ground height
	if birdCenterY+birdRadius >= groundY {
		return true
	}

	// Check if bird hit a pipe
	for _, pipe := range g.pipes {
		pipeLeft := float64(pipe.x)
		pipeRight := float64(pipe.x + pipeWidth)

		// Use a slightly larger collision radius to ensure we don't miss collisions
		// This helps prevent the bird from passing through pipes due to high speed
		extendedRadius := birdRadius * 1.2

		// Calculate horizontal distance between bird center and closest point on pipe
		closestX := math.Max(pipeLeft, math.Min(birdCenterX, pipeRight))
		horizontalDist := birdCenterX - closestX

		// If bird is not horizontally close to pipe, skip
		// Use the extended radius for this check
		if math.Abs(horizontalDist) > extendedRadius {
			continue
		}

		// Check collision with top pipe
		topPipeBottom := float64(pipe.gapStart)
		if birdCenterY-birdRadius < topPipeBottom {
			// Calculate vertical distance between bird center and bottom of top pipe
			verticalDist := birdCenterY - topPipeBottom

			// Check if bird is close enough to cause collision
			// Use the extended radius for this check
			if math.Abs(verticalDist) <= extendedRadius ||
				(math.Pow(horizontalDist, 2)+math.Pow(verticalDist, 2) <= math.Pow(extendedRadius, 2)) {
				return true
			}
		}

		// Check collision with bottom pipe
		bottomPipeTop := float64(pipe.gapStart + pipeGap)
		if birdCenterY+birdRadius > bottomPipeTop {
			// Calculate vertical distance between bird center and top of bottom pipe
			verticalDist := birdCenterY - bottomPipeTop

			// Check if bird is close enough to cause collision
			// Use the extended radius for this check
			if math.Abs(verticalDist) <= extendedRadius ||
				(math.Pow(horizontalDist, 2)+math.Pow(verticalDist, 2) <= math.Pow(extendedRadius, 2)) {
				return true
			}
		}
	}

	return false
}

// wouldCollideAfterJump simulates a jump and checks if it would lead to a collision
func (g *Game) wouldCollideAfterJump() bool {
	// Save current state
	originalY := g.bird.y
	originalVelocity := g.bird.velocity

	// Simulate jump and several frames ahead
	g.bird.velocity = jumpForce

	// Simulate 5 frames ahead (reduced from 10 to make AI less cautious)
	for i := 0; i < 5; i++ {
		g.bird.y += g.bird.velocity
		g.bird.velocity += gravity

		// Check collision at each frame
		if g.wouldCollide() {
			// Restore state and return true
			g.bird.y = originalY
			g.bird.velocity = originalVelocity
			return true
		}
	}

	// Restore state
	g.bird.y = originalY
	g.bird.velocity = originalVelocity

	return false
}

// wouldCollideWithoutJump simulates not jumping and checks if it would lead to a collision
func (g *Game) wouldCollideWithoutJump() bool {
	// Save current state
	originalY := g.bird.y
	originalVelocity := g.bird.velocity

	// Simulate gravity without jump for several frames ahead
	tempVelocity := g.bird.velocity

	// Simulate 5 frames ahead (reduced from 10 to make AI less cautious)
	for i := 0; i < 5; i++ {
		tempVelocity += gravity
		g.bird.y += tempVelocity

		// Check collision at each frame
		if g.wouldCollide() {
			// Restore state and return true
			g.bird.y = originalY
			g.bird.velocity = originalVelocity
			return true
		}
	}

	// Restore state
	g.bird.y = originalY
	g.bird.velocity = originalVelocity

	return false
}

// controlBirdAI implements perfect AI control for the bird
func (g *Game) controlBirdAI() {
	// Find the next pipe that the bird needs to navigate
	var nextPipe *Pipe
	var secondPipe *Pipe
	pipeCount := 0

	for _, pipe := range g.pipes {
		if float64(pipe.x+pipeWidth) >= g.bird.x {
			if pipeCount == 0 {
				nextPipe = &pipe
				pipeCount++
			} else if pipeCount == 1 {
				secondPipe = &pipe
				break
			}
		}
	}

	// Emergency collision avoidance - highest priority
	// Check if current position would lead to collision
	if g.wouldCollide() {
		// Already colliding, try to jump
		g.bird.velocity = jumpForce
		return
	}

	// Check if continuing without jumping would lead to collision in the next few frames
	if g.wouldCollideWithoutJump() {
		// Would collide if we don't jump, so jump
		g.bird.velocity = jumpForce
		return
	}

	// Check if jumping would lead to collision in the next few frames
	if g.wouldCollideAfterJump() {
		// Would collide if we jump, so don't jump
		return
	}

	// If we have a pipe to navigate
	if nextPipe != nil {
		// Calculate the center of the gap
		gapCenter := float64(nextPipe.gapStart) + float64(pipeGap)/2

		// Calculate bird's center
		birdCenterY := g.bird.y + float64(g.bird.height)/2

		// Calculate horizontal distance to the pipe
		distanceToPipe := float64(nextPipe.x) - g.bird.x

		// Calculate bird's vertical distance to gap center
		verticalDistToGap := birdCenterY - gapCenter

		// More accurate prediction of bird's position at the pipe
		// Calculate how many frames it will take to reach the pipe
		framesUntilPipe := distanceToPipe / pipeSpeed

		// Simulate bird's trajectory to predict position at pipe
		predictedY := birdCenterY
		predictedVelocity := g.bird.velocity

		// Simulate gravity effect over those frames
		for i := 0; i < int(framesUntilPipe); i++ {
			predictedVelocity += gravity
			predictedY += predictedVelocity
		}

		// Look ahead to the second pipe if it exists
		if secondPipe != nil && distanceToPipe < 250 { // Increased from 200 to look ahead earlier
			// Calculate the center of the second gap
			secondGapCenter := float64(secondPipe.gapStart) + float64(pipeGap)/2

			// If the second pipe's gap is significantly higher or lower, start adjusting early
			if secondGapCenter < gapCenter-40 || secondGapCenter > gapCenter+40 { // Reduced threshold from 50 to 40
				// Blend the target between current and next gap based on distance
				blendFactor := math.Max(0, math.Min(1, (250-distanceToPipe)/150)) // Adjusted blend calculation
				targetY := gapCenter*(1-blendFactor) + secondGapCenter*blendFactor

				// Adjust our target based on the blend
				gapCenter = targetY
			}
		}

		// Different strategies based on distance to pipe
		if distanceToPipe < 120 { // Increased from 100 to start precise positioning earlier
			// Very close to pipe - precise positioning is critical
			// Use the more accurate predicted position

			// Calculate safe zone (wider than before to make navigation easier)
			safeZoneTop := gapCenter - pipeGap/3
			safeZoneBottom := gapCenter + pipeGap/3

			// If we're going to be outside the safe zone, correct immediately
			if predictedY < safeZoneTop {
				// We'll be too high, let gravity work
				return
			} else if predictedY > safeZoneBottom {
				// We'll be too low, jump
				g.bird.velocity = jumpForce
				return
			}

			// Fine-tuning based on current position and velocity
			if verticalDistToGap < -20 && g.bird.velocity > 0 { // Increased threshold from -15 to -20
				// Above gap and moving down - let gravity work
				return
			} else if verticalDistToGap > 20 && g.bird.velocity < 0 { // Increased threshold from 15 to 20
				// Below gap and moving up - let momentum carry
				return
			} else if verticalDistToGap > 20 { // Increased threshold from 15 to 20
				// Below gap and not moving up enough - jump
				g.bird.velocity = jumpForce
				return
			}

			// Additional check for velocity - allow more speed variation
			if g.bird.velocity > 4 { // Increased from 3 to 4
				// Moving down too fast
				g.bird.velocity = jumpForce
				return
			} else if g.bird.velocity < -4 { // Increased from -3 to -4
				// Moving up too fast
				return
			}
		} else if distanceToPipe < 250 { // Increased from 200 to start aligning earlier
			// Medium distance - start aligning with gap
			// Use predicted position for better decision making
			verticalDistToPredicted := predictedY - gapCenter

			if verticalDistToPredicted < -35 { // Increased threshold from -25 to -35
				// Predicted to be too high, let gravity work
				return
			} else if verticalDistToPredicted > 35 { // Increased threshold from 25 to 35
				// Predicted to be too low, jump
				g.bird.velocity = jumpForce
				return
			} else if math.Abs(verticalDistToGap) < 25 && g.bird.velocity > 2.0 { // Increased threshold from 20 to 25, increased velocity threshold from 1.5 to 2.0
				// Near center but moving down too fast
				g.bird.velocity = jumpForce
				return
			}
		} else {
			// Far from pipe - maintain a good position
			// Use a more central target position
			targetY := float64(screenHeight) / 2

			// Blend between middle of screen and gap center based on distance
			blendFactor := math.Min(1, distanceToPipe/300) // Reduced from 400 to 300 to focus more on gap center
			targetY = gapCenter*(1-blendFactor) + targetY*blendFactor

			if birdCenterY > targetY+40 { // Increased from 30 to 40
				// Too low, jump
				g.bird.velocity = jumpForce
				return
			} else if birdCenterY < targetY-60 && g.bird.velocity < 0 { // Increased from 50 to 60
				// Too high and still going up
				return
			}
		}
	} else {
		// If no pipe is visible, maintain a middle position
		birdCenterY := g.bird.y + float64(g.bird.height)/2
		middleY := float64(screenHeight) / 2.2 // Slightly lower than center to prepare for pipes

		// Keep bird in the middle third of the screen with more flexibility
		if birdCenterY > middleY+40 { // Increased from 30 to 40
			// Too low, jump
			g.bird.velocity = jumpForce
		} else if birdCenterY < middleY-60 && g.bird.velocity < 0 { // Increased from 50 to 60
			// Too high and still going up
			return
		}

		// Don't let velocity get too extreme, but allow more variation
		if g.bird.velocity > 4 { // Increased from 3 to 4
			// Moving down too fast
			g.bird.velocity = jumpForce
		} else if g.bird.velocity < -4 { // Increased from -3 to -4
			// Moving up too fast
			return
		}
	}
}

// Draw draws the game state to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background
	screen.Fill(backgroundColor)

	// Draw common background elements
	g.drawBackground(screen)

	// Draw game elements
	g.drawGameElements(screen)

	// Draw state-specific UI
	switch g.state {
	case stateMenu:
		g.drawMenu(screen)
	case statePlay:
		g.drawGameUI(screen)
	}
}

// drawBackground draws background elements like clouds and ground
func (g *Game) drawBackground(screen *ebiten.Image) {
	// Draw background image
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.backgroundImg, op)

	// Draw clouds
	for _, cloud := range g.clouds {
		// Draw main cloud body
		ebitenutil.DrawCircle(screen, cloud.x, cloud.y, cloud.width/2, cloudColor)
		ebitenutil.DrawCircle(screen, cloud.x+cloud.width/3, cloud.y-cloud.height/4, cloud.width/3, cloudColor)
		ebitenutil.DrawCircle(screen, cloud.x+cloud.width/2, cloud.y, cloud.width/3, cloudColor)
	}

	// Draw ground
	groundHeight := 30.0
	groundY := float64(screenHeight) - groundHeight

	// Draw ground image with scrolling
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.groundOffset, groundY)
	screen.DrawImage(g.groundImg, op)

	// Draw a second ground image to create a seamless scrolling effect
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.groundOffset+float64(screenWidth), groundY)
	screen.DrawImage(g.groundImg, op)
}

// drawGameElements draws the bird and pipes
func (g *Game) drawGameElements(screen *ebiten.Image) {
	// Draw bird
	op := &ebiten.DrawImageOptions{}

	// Apply rotation based on velocity for more natural movement
	rotationAngle := math.Atan2(g.bird.velocity, pipeSpeed) // Use velocity to determine angle

	// Rotate around the center of the bird
	op.GeoM.Translate(-float64(g.bird.width)/2, -float64(g.bird.height)/2) // Move to center
	op.GeoM.Rotate(rotationAngle)                                          // Apply rotation
	op.GeoM.Translate(float64(g.bird.width)/2, float64(g.bird.height)/2)   // Move back

	// Set bird position
	op.GeoM.Translate(g.bird.x, g.bird.y)

	// Draw the bird image
	screen.DrawImage(g.birdImg, op)

	// Draw pipes
	for _, pipe := range g.pipes {
		pipeX := float64(pipe.x)

		// Draw top pipe (upside down)
		topPipeHeight := float64(pipe.gapStart)

		// Pipe body
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, topPipeHeight/100) // Scale to correct height (pipe image is 100px tall)
		op.GeoM.Translate(pipeX, 0)
		screen.DrawImage(g.pipeImg, op)

		// Pipe cap
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pipeX-5, topPipeHeight-20) // Position cap at bottom of top pipe
		screen.DrawImage(g.pipeCapImg, op)

		// Draw bottom pipe
		bottomPipeY := float64(pipe.gapStart + pipeGap)
		bottomPipeHeight := float64(screenHeight - pipe.gapStart - pipeGap)

		// Pipe body
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, bottomPipeHeight/100) // Scale to correct height
		op.GeoM.Translate(pipeX, bottomPipeY)
		screen.DrawImage(g.pipeImg, op)

		// Pipe cap
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pipeX-5, bottomPipeY) // Position cap at top of bottom pipe
		screen.DrawImage(g.pipeCapImg, op)
	}
}

// drawMenu draws the main menu UI
func (g *Game) drawMenu(screen *ebiten.Image) {
	// Draw semi-transparent overlay to make menu text more visible
	overlayColor := color.RGBA{0, 0, 0, 128} // Semi-transparent black
	ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), overlayColor)

	// Draw game title
	titleText := "FLAPPY BIRD"
	ebitenutil.DebugPrintAt(screen, titleText, screenWidth/2-50, screenHeight/4)

	// Draw menu options
	startText := "Start Game"
	exitText := "Exit"

	// Highlight selected option
	if g.menuOption == 0 {
		startText = "> " + startText + " <"
	} else {
		exitText = "> " + exitText + " <"
	}

	ebitenutil.DebugPrintAt(screen, startText, screenWidth/2-50, screenHeight/2)
	ebitenutil.DebugPrintAt(screen, exitText, screenWidth/2-50, screenHeight/2+30)

	// Draw controls help
	controlsText := "Use UP/DOWN to navigate, SPACE/ENTER to select"
	ebitenutil.DebugPrintAt(screen, controlsText, screenWidth/2-150, screenHeight*3/4)

	// Draw demo indicator
	demoText := "DEMO PLAYING"
	ebitenutil.DebugPrintAt(screen, demoText, screenWidth-100, 30)

	// Draw score in demo
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 30)
}

// drawGameUI draws the gameplay UI
func (g *Game) drawGameUI(screen *ebiten.Image) {
	// Draw score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 30)

	// Draw game over message
	if g.isGameOver {
		gameOverText := "Game Over! Press 'R' to restart."
		ebitenutil.DebugPrintAt(screen, gameOverText, screenWidth/2-150, screenHeight/2)
	}

	// Draw controls help
	controlsText := "Press ESC to return to menu"
	ebitenutil.DebugPrintAt(screen, controlsText, screenWidth-200, screenHeight-20)
}

// Layout implements ebiten.Game's Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// checkCollisions checks if the bird has collided with pipes or boundaries
func (g *Game) checkCollisions() {
	// Calculate bird center and radius for circle-based collision
	birdCenterX := g.bird.x + float64(g.bird.width)/2
	birdCenterY := g.bird.y + float64(g.bird.height)/2
	birdRadius := float64(g.bird.width) / 2

	// Check if bird hit the top boundary
	if birdCenterY-birdRadius <= 0 {
		g.isGameOver = true
		// Play collision sound
		g.collisionSound.Rewind()
		g.collisionSound.Play()
		return
	}

	// Check if bird hit the ground
	groundY := float64(screenHeight) - 30.0 // 30 is the ground height
	if birdCenterY+birdRadius >= groundY {
		g.isGameOver = true
		// Play collision sound
		g.collisionSound.Rewind()
		g.collisionSound.Play()
		return
	}

	// Check if bird hit a pipe
	for _, pipe := range g.pipes {
		pipeLeft := float64(pipe.x)
		pipeRight := float64(pipe.x + pipeWidth)

		// Use a slightly larger collision radius to ensure we don't miss collisions
		// This helps prevent the bird from passing through pipes due to high speed
		extendedRadius := birdRadius * 1.2

		// Calculate horizontal distance between bird center and closest point on pipe
		closestX := math.Max(pipeLeft, math.Min(birdCenterX, pipeRight))
		horizontalDist := birdCenterX - closestX

		// If bird is not horizontally close to pipe, skip
		// Use the extended radius for this check
		if math.Abs(horizontalDist) > extendedRadius {
			continue
		}

		// Check collision with top pipe
		topPipeBottom := float64(pipe.gapStart)
		if birdCenterY-birdRadius < topPipeBottom {
			// Calculate vertical distance between bird center and bottom of top pipe
			verticalDist := birdCenterY - topPipeBottom

			// Check if bird is close enough to cause collision
			// Use the extended radius for this check
			if math.Abs(verticalDist) <= extendedRadius ||
				(math.Pow(horizontalDist, 2)+math.Pow(verticalDist, 2) <= math.Pow(extendedRadius, 2)) {
				g.isGameOver = true
				// Play collision sound
				g.collisionSound.Rewind()
				g.collisionSound.Play()
				return
			}
		}

		// Check collision with bottom pipe
		bottomPipeTop := float64(pipe.gapStart + pipeGap)
		if birdCenterY+birdRadius > bottomPipeTop {
			// Calculate vertical distance between bird center and top of bottom pipe
			verticalDist := birdCenterY - bottomPipeTop

			// Check if bird is close enough to cause collision
			// Use the extended radius for this check
			if math.Abs(verticalDist) <= extendedRadius ||
				(math.Pow(horizontalDist, 2)+math.Pow(verticalDist, 2) <= math.Pow(extendedRadius, 2)) {
				g.isGameOver = true
				// Play collision sound
				g.collisionSound.Rewind()
				g.collisionSound.Play()
				return
			}
		}
	}
}

// Reset resets the game state
func (g *Game) Reset() {
	// Save current state and demo mode
	currentState := g.state
	isDemoMode := g.demoMode

	g.bird = Bird{
		x:         float64(screenWidth) / 4,
		y:         float64(screenHeight) / 2,
		width:     birdSize,
		height:    birdSize,
		wingAngle: 0,
	}
	g.pipes = []Pipe{}

	// Reset clouds to new random positions
	for i := range g.clouds {
		g.clouds[i] = Cloud{
			x:      float64(rand.Intn(screenWidth)),
			y:      float64(rand.Intn(screenHeight / 3)),
			width:  float64(rand.Intn(80) + 40),
			height: float64(rand.Intn(30) + 20),
			speed:  float64(rand.Intn(2) + 1),
		}
	}

	g.groundOffset = 0
	g.score = 0
	g.isGameOver = false
	g.tick = 0

	// Restore state and demo mode
	g.state = currentState
	g.demoMode = isDemoMode

	// Clear key states
	g.keyStates = make(map[ebiten.Key]bool)
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Set up the game
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Bird")

	// Create and run the game
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
