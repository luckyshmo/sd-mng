package storage

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const originalPath = "original"
const upscaledPath = "upscaled"

type MangaInfo struct {
	Title             string
	Cover             []byte
	VolumeCount       int
	IsUpscaled        bool
	AverageMegapixels float64
}

type Storage interface {
	GetStoredMangaInfo() ([]*MangaInfo, error)
	NewReader(title string) Reader
	Write(images []string, names []string, title, relativePath string) error
}

type FSStorage struct {
	upscaledRoot string
	originalRoot string
}

func init() {
	rootDir := os.Getenv("MANGA_STORAGE_DIR")
	err1 := os.MkdirAll(rootDir+"/"+originalPath, 0755)
	err2 := os.MkdirAll(rootDir+"/"+upscaledPath, 0755)
	if err1 != nil || err2 != nil {
		log.Fatalf("failed to create storage directories: %v, %v", err1, err2)
	}
}

func NewFSStorage() *FSStorage {
	rootDir := os.Getenv("MANGA_STORAGE_DIR")
	return &FSStorage{
		upscaledRoot: rootDir + "/" + upscaledPath,
		originalRoot: rootDir + "/" + originalPath,
	}
}

func (s *FSStorage) NewReader(title string) Reader {
	return NewFSReader(s.originalRoot + "/" + title)
}

func (s *FSStorage) GetStoredMangaInfo() ([]*MangaInfo, error) {
	entities, err := os.ReadDir(s.originalRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to read original manga directory: %w", err)
	}

	var mangaInfo []*MangaInfo
	for _, entity := range entities {
		if !entity.IsDir() {
			continue
		}

		volumeEntities, err := os.ReadDir(s.originalRoot + "/" + entity.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read volume directory: %w", err)
		}

		volumeCount := 0
		var cover []byte
		for _, volumeEntity := range volumeEntities {
			if volumeEntity.IsDir() {
				volumeCount++
				if cover == nil {
					cover, err = os.ReadFile(s.originalRoot + "/" + entity.Name() + "/" + volumeEntity.Name() + "/0.jpeg")
					if err != nil {
						log.Printf("failed to read file %s: %v", volumeEntity.Name(), err)
						cover = nil
						continue
					}
				}
			}
		}

		mangaInfo = append(mangaInfo, &MangaInfo{
			Title:       entity.Name(),
			Cover:       cover,
			VolumeCount: volumeCount,
		})
	}

	return mangaInfo, nil
}

func (s *FSStorage) Write(images []string, names []string, title, relativePath string) error {
	if _, err := os.Stat(s.upscaledRoot + "/" + title + "/" + relativePath); os.IsNotExist(err) {
		err := os.MkdirAll(s.upscaledRoot+"/"+title+"/"+relativePath, 0755)
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
		img, err := convertToJPEG(imageData)
		if err != nil {
			return fmt.Errorf("convert to JPEG: %w", err)
		}
		elapsed += int64(time.Since(now))

		imagePath := filepath.Join(s.upscaledRoot, title, relativePath, names[i])
		err = os.WriteFile(imagePath, img, 0644)
		if err != nil {
			return err
		}
	}

	fmt.Println("average time spent is: ", elapsed/int64(len(images)), " ns")

	return nil
}
