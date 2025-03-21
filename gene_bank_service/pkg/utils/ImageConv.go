package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// fetchImageFromURL fetches the image as bytes from the provided URL
func FetchImageFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the image content into a byte slice
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return imageBytes, nil
}

// fetchImageFromPath reads the image from a local file path and returns it as a byte slice
func FetchImageFromPath(imagePath string) ([]byte, error) {
	var imageBytes []byte
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image from path: %v", err)
	}

	return imageBytes, nil
}
