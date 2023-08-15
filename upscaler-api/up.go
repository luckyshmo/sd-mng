package upscale

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	APIURL         = "http://192.168.1.10:7860/sdapi/v1/extra-batch-images"
	OriginFolder   = "./origin"
	UpscaledFolder = "./upscaled"

	Upscaler = "R-ESRGAN 4x+ Anime6B"
)

type APIResponse struct {
	HTMLInfo string   `json:"html_info"`
	Images   []string `json:"images"`
}

func Upscale(source, dist string) {
	OriginFolder = source
	UpscaledFolder = dist

	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	ch := make(chan ImageInfo, 1024)
	go reader(OriginFolder, "", ch, cancelFunc)
	go writer(ch, ctx)

	<-ctx.Done()
	fmt.Println("Done")
	time.Sleep(time.Hour * 24)
}

func reader(basePath, relativePath string, ch chan ImageInfo, finish context.CancelFunc) {
	defer finish()
	processFolder(basePath, relativePath, ch)
}

func writer(ch chan ImageInfo, ctx context.Context) {
	for {
		select {
		case info := <-ch:
			fmt.Println("process: ", info.path)
			processImages(info)
			fmt.Println("finish: ", info.path)
		}
	}
}

func processImages(info ImageInfo) {
	// Create the request payload
	payload := map[string]interface{}{
		"resize_mode":                  0,
		"show_extras_results":          true,
		"gfpgan_visibility":            0,
		"codeformer_visibility":        0,
		"codeformer_weight":            0,
		"upscaling_resize":             2,
		"upscaling_resize_w":           512,
		"upscaling_resize_h":           512,
		"upscaling_crop":               true,
		"upscaler_1":                   Upscaler,
		"upscaler_2":                   "None",
		"extras_upscaler_2_visibility": 0,
		"upscale_first":                false,
		"imageList":                    info.imageList,
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling payload: %v\n", err)
		return
	}

	// Send POST request to the API
	resp, err := http.Post(APIURL, "application/json", bytes.NewBuffer(jsonPayload))
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
		fmt.Println("empty image response for ", info.path)
		return
	}

	// Store upscaled images in the upscaled folder
	err = storeImages(apiResponse.Images, info.names, info.path)
	if err != nil {
		fmt.Printf("Error storing upscaled images: %v\n", err)
		return
	}
}

type ImageInfo struct {
	imageList []map[string]string
	names     []string
	path      string
}

func processFolder(basePath, relativePath string, ch chan ImageInfo) {
	path := filepath.Join(basePath, relativePath)
	fmt.Println("processing path: ", path)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("ERROR: read dir %s: %v\n", path, err)
		return
	}

	var imageList []map[string]string
	names := make([]string, len(files))

	for i, f := range files {
		if f.IsDir() {
			processFolder(basePath, relativePath+"/"+f.Name(), ch)
			continue
		}

		names[i] = f.Name()

		imageData, err := os.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			log.Printf("ERROR: reading image '%s': %s", filepath.Join(path, f.Name()), err.Error())
			continue
		}

		imageList = append(imageList, map[string]string{
			"data": imageDataBase64(imageData),
			"name": f.Name(),
		})
	}

	if len(imageList) > 0 {
		ch <- ImageInfo{
			imageList: imageList,
			names:     names,
			path:      relativePath,
		}
	}
}

// imageDataBase64 encodes image data as base64.
func imageDataBase64(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}

// storeImages stores the upscaled images in a folder.
func storeImages(images []string, names []string, path string) error {
	if _, err := os.Stat(UpscaledFolder + "/" + path); os.IsNotExist(err) {
		err := os.MkdirAll(UpscaledFolder+"/"+path, 0755)
		if err != nil {
			return err
		}
	}

	var elapsed int64 = 0
	for i, image := range images {
		imageData, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			return err
		}

		now := time.Now()
		img, err := ConvertToJPEG(imageData)
		if err != nil {
			return fmt.Errorf("convert to JPEG: %w", err)
		}
		elapsed += int64(time.Since(now))

		imagePath := filepath.Join(UpscaledFolder, path, names[i])
		err = os.WriteFile(imagePath, img, 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("average time spent is: ", elapsed/int64(len(images)), " ns")

	return nil
}
