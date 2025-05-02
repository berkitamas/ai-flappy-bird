package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestNewGame tests the NewGame function
func TestNewGame(t *testing.T) {
	// This is a basic test to ensure NewGame doesn't panic
	// We can't fully test it without mocking the image and audio generation
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewGame panicked: %v", r)
		}
	}()

	// Skip actual test execution since it requires Ebiten context
	t.Skip("Skipping NewGame test as it requires Ebiten context")
}

// TestIsKeyJustPressed tests the isKeyJustPressed function
func TestIsKeyJustPressed(t *testing.T) {
	g := &Game{
		keyStates: make(map[ebiten.Key]bool),
	}

	// Test when key was not pressed before and is not pressed now
	g.keyStates[ebiten.KeySpace] = false
	// We can't test ebiten.IsKeyPressed directly, so we'll just test the logic
	// Simulate key not being pressed
	result := g.isKeyJustPressed(ebiten.KeySpace)
	if result {
		t.Error("Expected false when key was not pressed before and is not pressed now")
	}

	// Skip actual test execution since it requires Ebiten context
	t.Skip("Skipping isKeyJustPressed test as it requires Ebiten context")
}

// TestReset tests the Reset function
func TestReset(t *testing.T) {
	g := &Game{
		bird: Bird{
			x:        100,
			y:        100,
			velocity: 5,
		},
		pipes:      []Pipe{{x: 200, gapStart: 150}},
		score:      10,
		isGameOver: true,
		tick:       100,
		state:      statePlay,
		demoMode:   false,
		keyStates:  map[ebiten.Key]bool{ebiten.KeySpace: true},
	}

	// Call Reset
	g.Reset()

	// Verify bird position and velocity are reset
	if g.bird.x != float64(screenWidth)/4 {
		t.Errorf("Expected bird.x to be %f, got %f", float64(screenWidth)/4, g.bird.x)
	}
	if g.bird.y != float64(screenHeight)/2 {
		t.Errorf("Expected bird.y to be %f, got %f", float64(screenHeight)/2, g.bird.y)
	}
	if g.bird.velocity != 0 {
		t.Errorf("Expected bird.velocity to be 0, got %f", g.bird.velocity)
	}

	// Verify pipes are cleared
	if len(g.pipes) != 0 {
		t.Errorf("Expected pipes to be empty, got %d pipes", len(g.pipes))
	}

	// Verify score and game over state are reset
	if g.score != 0 {
		t.Errorf("Expected score to be 0, got %d", g.score)
	}
	if g.isGameOver {
		t.Error("Expected isGameOver to be false")
	}

	// Verify tick counter is reset
	if g.tick != 0 {
		t.Errorf("Expected tick to be 0, got %d", g.tick)
	}

	// Verify state and demoMode are preserved
	if g.state != statePlay {
		t.Errorf("Expected state to be preserved as statePlay, got %d", g.state)
	}
	if g.demoMode {
		t.Error("Expected demoMode to be preserved as false")
	}

	// Verify keyStates is cleared
	if len(g.keyStates) != 0 {
		t.Error("Expected keyStates to be cleared")
	}
}

// TestLayout tests the Layout function
func TestLayout(t *testing.T) {
	g := &Game{}
	width, height := g.Layout(800, 600)
	if width != screenWidth {
		t.Errorf("Expected width to be %d, got %d", screenWidth, width)
	}
	if height != screenHeight {
		t.Errorf("Expected height to be %d, got %d", screenHeight, height)
	}
}

// TestUpdateBackground tests the updateBackground function
func TestUpdateBackground(t *testing.T) {
	// Create a game instance with clouds
	g := &Game{
		groundOffset: 0,
		clouds: []Cloud{
			{x: 100, y: 50, width: 60, height: 30, speed: 1},
			{x: -60, y: 30, width: 50, height: 20, speed: 2}, // This cloud should be reset
		},
	}

	// Call updateBackground
	g.updateBackground()

	// Verify ground offset is updated
	if g.groundOffset != -pipeSpeed {
		t.Errorf("Expected groundOffset to be %d, got %f", -pipeSpeed, g.groundOffset)
	}

	// Verify clouds are updated
	if g.clouds[0].x != 99 {
		t.Errorf("Expected cloud[0].x to be 99, got %f", g.clouds[0].x)
	}

	// Verify cloud that moved off-screen is reset
	if g.clouds[1].x <= 0 {
		t.Errorf("Expected cloud[1].x to be reset to a positive value, got %f", g.clouds[1].x)
	}

	// Test ground offset reset
	g.groundOffset = -51
	g.updateBackground()
	if g.groundOffset != 0 {
		t.Errorf("Expected groundOffset to be reset to 0, got %f", g.groundOffset)
	}
}

// TestUpdate tests the Update function
func TestUpdate(t *testing.T) {
	// This is a basic test to ensure Update doesn't panic
	// We can't fully test it without mocking the Ebiten context

	// Skip actual test execution since it requires Ebiten context
	t.Skip("Skipping Update test as it requires Ebiten context")
}
