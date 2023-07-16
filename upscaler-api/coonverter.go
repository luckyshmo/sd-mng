package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
)

// ConvertToJPEG converts the input image byte slice to JPEG format.
func ConvertToJPEG(img []byte) ([]byte, error) {
	// Create a buffer to store the JPEG image
	jpegBuffer := new(bytes.Buffer)

	decodedImg, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, fmt.Errorf("decode :%w", err)
	}

	// Encode the image as JPEG
	err = jpeg.Encode(jpegBuffer, decodedImg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image as JPEG: %v", err)
	}

	return jpegBuffer.Bytes(), nil
}
