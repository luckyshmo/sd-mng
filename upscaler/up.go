package upscaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kek.com/storage"
)

type Upscaler interface {
	Upscale(info *storage.ImageInfo) (images []string, names []string, path string)
}

type UpscalerConfig struct {
	APIURL       string
	UpscalerType string
}

// SDUpscaler is a stable diffusion upscaler.
type SDUpscaler struct {
	UpscalerConfig
}

func NewSDUpscaler(cfg UpscalerConfig) *SDUpscaler {
	return &SDUpscaler{
		cfg,
	}
}

type APIResponse struct {
	HTMLInfo string   `json:"html_info"`
	Images   []string `json:"images"`
}

func (u *SDUpscaler) Upscale(info *storage.ImageInfo) (images []string, names []string, path string) {
	// Create the request payload
	payload := map[string]interface{}{
		"resize_mode":                  0,
		"show_extras_results":          true,
		"gfpgan_visibility":            0,
		"codeformer_visibility":        0,
		"codeformer_weight":            0,
		"upscaling_resize":             1.5,
		"upscaling_resize_w":           512,
		"upscaling_resize_h":           512,
		"upscaling_crop":               true,
		"upscaler_1":                   u.UpscalerType,
		"upscaler_2":                   "None",
		"extras_upscaler_2_visibility": 0,
		"upscale_first":                false,
		"imageList":                    info.ImageList,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}

	// Send POST request to the API
	resp, err := http.Post(u.APIURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Parse API response
	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Error parsing API response: %v\n", err)
		return
	}

	if len(apiResponse.Images) < 1 {
		fmt.Println("empty image response for ", info.Path)
		return
	}

	return apiResponse.Images, info.Names, info.Path
}
