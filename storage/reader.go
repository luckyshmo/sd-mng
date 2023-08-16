package storage

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ImageInfo struct {
	ImageList []map[string]string
	Names     []string
	Path      string
}

type Reader interface {
	Next() (*ImageInfo, bool)
}

type FSReader struct {
	img chan *ImageInfo
}

func NewFSReader(path string) *FSReader {
	reader := &FSReader{make(chan *ImageInfo, 1)}
	go func() {
		reader.read(path, "")
		close(reader.img)
	}()
	return reader
}

func (r *FSReader) Next() (*ImageInfo, bool) {
	img, ok := <-r.img
	return img, ok
}

func (r *FSReader) read(basePath, relativePath string) {
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
			// put processing of next folder on stack
			defer r.read(basePath, relativePath+"/"+f.Name())
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
		r.img <- &ImageInfo{
			ImageList: imageList,
			Names:     names,
			Path:      relativePath,
		}
	}
}

// imageDataBase64 encodes image data as base64.
func imageDataBase64(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}
