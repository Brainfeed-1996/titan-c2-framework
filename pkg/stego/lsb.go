/*
 * Package: stego
 * Author: Olivier Robert-Duboille
 * Description: Implementation of Least Significant Bit (LSB) image steganography for covert data transport.
 */

package stego

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
)

// Encoder handles embedding data into images
type Encoder struct{}

// Decoder handles extracting data from images
type Decoder struct{}

// NewEncoder creates a new steganography encoder
func NewEncoder() *Encoder {
	return &Encoder{}
}

// NewDecoder creates a new steganography decoder
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Embed hides the message bytes into the source image and writes result to writer.
// Uses Least Significant Bit (LSB) encoding.
// Format: [Length: 4 bytes][Message: N bytes]
func (e *Encoder) Embed(src image.Image, message []byte, w io.Writer) error {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Prepare data: length prefix + message
	length := uint32(len(message))
	data := make([]byte, 4+len(message))
	data[0] = byte(length >> 24)
	data[1] = byte(length >> 16)
	data[2] = byte(length >> 8)
	data[3] = byte(length)
	copy(data[4:], message)

	// Check capacity (3 bits per pixel: R, G, B)
	capacity := (width * height * 3) / 8
	if len(data) > capacity {
		return errors.New("image too small to hold message")
	}

	// Create mutable image
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	dataIdx := 0
	bitIdx := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if dataIdx >= len(data) {
				break
			}

			r, g, b, a := dst.At(x, y).RGBA()
			// Convert to 8-bit
			r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

			// Encode 3 bits per pixel (1 in R, 1 in G, 1 in B)
			
			// Red
			if dataIdx < len(data) {
				bit := (data[dataIdx] >> (7 - bitIdx)) & 1
				r8 = (r8 & 0xFE) | bit
				bitIdx++
				if bitIdx == 8 {
					bitIdx = 0
					dataIdx++
				}
			}

			// Green
			if dataIdx < len(data) {
				bit := (data[dataIdx] >> (7 - bitIdx)) & 1
				g8 = (g8 & 0xFE) | bit
				bitIdx++
				if bitIdx == 8 {
					bitIdx = 0
					dataIdx++
				}
			}

			// Blue
			if dataIdx < len(data) {
				bit := (data[dataIdx] >> (7 - bitIdx)) & 1
				b8 = (b8 & 0xFE) | bit
				bitIdx++
				if bitIdx == 8 {
					bitIdx = 0
					dataIdx++
				}
			}

			dst.Set(x, y, color.RGBA{r8, g8, b8, a8})
		}
	}

	return png.Encode(w, dst)
}

// Extract retrieves hidden data from an image
func (d *Decoder) Extract(src image.Image) ([]byte, error) {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var length uint32
	var message []byte
	
	// Reading state
	var currentByte byte
	bitIdx := 0
	
	// Phase 0: Reading Length (4 bytes)
	// Phase 1: Reading Message
	bytesRead := 0
	readLength := false

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			vals := []uint8{r8, g8, b8}

			for _, val := range vals {
				bit := val & 1
				currentByte = (currentByte << 1) | bit
				bitIdx++

				if bitIdx == 8 {
					if !readLength {
						// Accumulate length bytes
						if bytesRead < 4 {
							length = (length << 8) | uint32(currentByte)
							bytesRead++
						}
						if bytesRead == 4 {
							readLength = true
							message = make([]byte, 0, length)
							bytesRead = 0 // Reset for message counting
						}
					} else {
						// Append message byte
						message = append(message, currentByte)
						if uint32(len(message)) == length {
							return message, nil
						}
					}
					currentByte = 0
					bitIdx = 0
				}
			}
		}
	}

	return nil, errors.New("unexpected end of image data")
}
