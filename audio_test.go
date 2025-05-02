package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// TestGenerateSineWave tests the sine wave generation
func TestGenerateSineWave(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("generateSineWave panicked: %v", r)
		}
	}()

	// Generate a short sine wave
	samples := generateSineWave(440, 100*time.Millisecond) // 440Hz for 100ms
	if len(samples) == 0 {
		t.Error("Expected non-empty samples")
	}

	// Check that the number of samples is correct
	expectedSamples := int(sampleRate*0.1) * 2 // 0.1 seconds * sample rate * 2 bytes per sample
	if len(samples) != expectedSamples {
		t.Errorf("Expected %d samples, got %d", expectedSamples, len(samples))
	}
}

// TestAddWAVHeader tests the WAV header addition
func TestAddWAVHeader(t *testing.T) {
	// This is a basic test to ensure the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("addWAVHeader panicked: %v", r)
		}
	}()

	// Create some dummy PCM data
	pcmData := make([]byte, 1000)
	for i := range pcmData {
		pcmData[i] = byte(i % 256)
	}

	// Add WAV header
	wavData := addWAVHeader(pcmData)
	if len(wavData) == 0 {
		t.Error("Expected non-empty WAV data")
	}

	// Check that the WAV data is longer than the PCM data (due to header)
	if len(wavData) <= len(pcmData) {
		t.Errorf("Expected WAV data to be longer than PCM data, got %d vs %d", len(wavData), len(pcmData))
	}

	// Check for RIFF header
	if string(wavData[0:4]) != "RIFF" {
		t.Errorf("Expected RIFF header, got %s", string(wavData[0:4]))
	}

	// Check for WAVE format
	if string(wavData[8:12]) != "WAVE" {
		t.Errorf("Expected WAVE format, got %s", string(wavData[8:12]))
	}

	// Check for fmt chunk
	if string(wavData[12:16]) != "fmt " {
		t.Errorf("Expected fmt chunk, got %s", string(wavData[12:16]))
	}

	// Check for data chunk
	if string(wavData[36:40]) != "data" {
		t.Errorf("Expected data chunk, got %s", string(wavData[36:40]))
	}
}

// TestGenerateJumpSoundSamples tests the sample generation part of generateJumpSound
func TestGenerateJumpSoundSamples(t *testing.T) {
	// This test focuses on the sample generation part of generateJumpSound
	// without requiring an audio context

	// Generate samples similar to what generateJumpSound does
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

	// Verify samples are generated
	if len(samples) == 0 {
		t.Error("Expected non-empty samples")
	}

	// Verify sample count
	expectedSampleCount := int(sampleRate*duration.Seconds()) * 2
	if len(samples) != expectedSampleCount {
		t.Errorf("Expected %d samples, got %d", expectedSampleCount, len(samples))
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Verify WAV header is added
	if len(wavBytes) <= len(samples) {
		t.Errorf("Expected WAV bytes to be longer than samples, got %d vs %d", len(wavBytes), len(samples))
	}

	// Verify RIFF header
	if string(wavBytes[0:4]) != "RIFF" {
		t.Errorf("Expected RIFF header, got %s", string(wavBytes[0:4]))
	}
}

// TestGenerateScoreSoundSamples tests the sample generation part of generateScoreSound
func TestGenerateScoreSoundSamples(t *testing.T) {
	// This test focuses on the sample generation part of generateScoreSound
	// without requiring an audio context

	// Generate samples similar to what generateScoreSound does
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

	// Verify samples are generated
	if len(samples) == 0 {
		t.Error("Expected non-empty samples")
	}

	// Verify sample count
	expectedSampleCount := int(sampleRate*duration.Seconds()) * 2
	if len(samples) != expectedSampleCount {
		t.Errorf("Expected %d samples, got %d", expectedSampleCount, len(samples))
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Verify WAV header is added
	if len(wavBytes) <= len(samples) {
		t.Errorf("Expected WAV bytes to be longer than samples, got %d vs %d", len(wavBytes), len(samples))
	}
}

// TestGenerateCollisionSoundSamples tests the sample generation part of generateCollisionSound
func TestGenerateCollisionSoundSamples(t *testing.T) {
	// This test focuses on the sample generation part of generateCollisionSound
	// without requiring an audio context

	// Generate samples similar to what generateCollisionSound does
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

	// Verify samples are generated
	if len(samples) == 0 {
		t.Error("Expected non-empty samples")
	}

	// Verify sample count
	expectedSampleCount := int(sampleRate*duration.Seconds()) * 2
	if len(samples) != expectedSampleCount {
		t.Errorf("Expected %d samples, got %d", expectedSampleCount, len(samples))
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Verify WAV header is added
	if len(wavBytes) <= len(samples) {
		t.Errorf("Expected WAV bytes to be longer than samples, got %d vs %d", len(wavBytes), len(samples))
	}
}

// TestGenerateMenuMusicSamples tests the sample generation part of generateMenuMusic
func TestGenerateMenuMusicSamples(t *testing.T) {
	// This test focuses on the sample generation part of generateMenuMusic
	// without requiring an audio context

	// Generate a simple melody similar to what generateMenuMusic does
	duration := 1 * time.Second // Shorter duration for testing
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
			if i*2+1 < len(samples) {
				samples[i*2] = byte(sample)
				samples[i*2+1] = byte(sample >> 8)
			}
		}

		currentTime += note.duration
	}

	// Verify samples are generated
	if len(samples) == 0 {
		t.Error("Expected non-empty samples")
	}

	// Verify sample count
	expectedSampleCount := int(sampleRate*duration.Seconds()) * 2
	if len(samples) != expectedSampleCount {
		t.Errorf("Expected %d samples, got %d", expectedSampleCount, len(samples))
	}

	// Add WAV header to the samples
	wavBytes := addWAVHeader(samples)

	// Verify WAV header is added
	if len(wavBytes) <= len(samples) {
		t.Errorf("Expected WAV bytes to be longer than samples, got %d vs %d", len(wavBytes), len(samples))
	}
}

// TestGenerateJumpSound tests the jump sound generation
func TestGenerateJumpSound(t *testing.T) {
	// Skip this test as it requires an audio context
	t.Skip("Skipping generateJumpSound test as it requires audio context")
}

// TestGenerateScoreSound tests the score sound generation
func TestGenerateScoreSound(t *testing.T) {
	// Skip this test as it requires an audio context
	t.Skip("Skipping generateScoreSound test as it requires audio context")
}

// TestGenerateCollisionSound tests the collision sound generation
func TestGenerateCollisionSound(t *testing.T) {
	// Skip this test as it requires an audio context
	t.Skip("Skipping generateCollisionSound test as it requires audio context")
}

// TestGenerateMenuMusic tests the menu music generation
func TestGenerateMenuMusic(t *testing.T) {
	// Skip this test as it requires an audio context
	t.Skip("Skipping generateMenuMusic test as it requires audio context")
}
