/*
 * Package: stego
 * Author: Olivier Robert-Duboille
 * Description: Unit tests for LSB steganography implementation.
 */

package stego

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func createTestImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 100, 255})
		}
	}
	return img
}

func TestEmbedExtract(t *testing.T) {
	// Create a small random image
	srcImg := createTestImage(100, 100)
	
	secretMessage := []byte("This is a covert C2 instruction.")
	
	// Embed
	enc := NewEncoder()
	var buf bytes.Buffer
	err := enc.Embed(srcImg, secretMessage, &buf)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}

	// Decode buffer back to image
	encodedImg, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("Failed to decode PNG: %v", err)
	}

	// Extract
	dec := NewDecoder()
	extracted, err := dec.Extract(encodedImg)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if string(extracted) != string(secretMessage) {
		t.Errorf("Expected '%s', got '%s'", secretMessage, extracted)
	}
}
