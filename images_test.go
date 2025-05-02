package main

import (
	"testing"
)

// TestGenerateBirdImage tests the bird image generation
func TestGenerateBirdImage(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generateBirdImage panicked: %v", r)
		}
	}()

	img := generateBirdImage()
	if img == nil {
		t.Error("Expected non-nil image")
	}

	// Check image dimensions
	width, height := img.Size()
	if width != birdSize {
		t.Errorf("Expected width to be %d, got %d", birdSize, width)
	}
	if height != birdSize {
		t.Errorf("Expected height to be %d, got %d", birdSize, height)
	}
}

// TestGeneratePipeImage tests the pipe image generation
func TestGeneratePipeImage(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generatePipeImage panicked: %v", r)
		}
	}()

	img := generatePipeImage()
	if img == nil {
		t.Error("Expected non-nil image")
	}

	// Check image dimensions
	width, height := img.Size()
	if width != pipeWidth {
		t.Errorf("Expected width to be %d, got %d", pipeWidth, width)
	}
	if height != 100 {
		t.Errorf("Expected height to be 100, got %d", height)
	}
}

// TestGeneratePipeCapImage tests the pipe cap image generation
func TestGeneratePipeCapImage(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generatePipeCapImage panicked: %v", r)
		}
	}()

	img := generatePipeCapImage()
	if img == nil {
		t.Error("Expected non-nil image")
	}

	// Check image dimensions
	width, height := img.Size()
	if width != pipeWidth+10 {
		t.Errorf("Expected width to be %d, got %d", pipeWidth+10, width)
	}
	if height != 20 {
		t.Errorf("Expected height to be 20, got %d", height)
	}
}

// TestGenerateBackgroundImage tests the background image generation
func TestGenerateBackgroundImage(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generateBackgroundImage panicked: %v", r)
		}
	}()

	img := generateBackgroundImage()
	if img == nil {
		t.Error("Expected non-nil image")
	}

	// Check image dimensions
	width, height := img.Size()
	if width != screenWidth {
		t.Errorf("Expected width to be %d, got %d", screenWidth, width)
	}
	if height != screenHeight {
		t.Errorf("Expected height to be %d, got %d", screenHeight, height)
	}
}

// TestGenerateGroundImage tests the ground image generation
func TestGenerateGroundImage(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generateGroundImage panicked: %v", r)
		}
	}()

	img := generateGroundImage()
	if img == nil {
		t.Error("Expected non-nil image")
	}

	// Check image dimensions
	width, height := img.Size()
	if width != screenWidth {
		t.Errorf("Expected width to be %d, got %d", screenWidth, width)
	}
	if height != 30 {
		t.Errorf("Expected height to be 30, got %d", height)
	}
}
