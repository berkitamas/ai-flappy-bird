package main

import (
	"bytes"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

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
