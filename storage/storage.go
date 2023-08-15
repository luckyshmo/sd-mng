package storage

import (
	"fmt"
	"log"
	"os"
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
}

type FSStorage struct {
	upscaledPath string
	originalPath string
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
		upscaledPath: rootDir + "/" + upscaledPath,
		originalPath: rootDir + "/" + originalPath,
	}
}

func (s *FSStorage) GetStoredMangaInfo() ([]*MangaInfo, error) {
	entities, err := os.ReadDir(s.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read original manga directory: %w", err)
	}

	var mangaInfo []*MangaInfo
	for _, entity := range entities {
		if !entity.IsDir() {
			continue
		}

		volumeEntities, err := os.ReadDir(s.originalPath + "/" + entity.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read volume directory: %w", err)
		}

		volumeCount := 0
		var cover []byte
		for _, volumeEntity := range volumeEntities {
			if volumeEntity.IsDir() {
				volumeCount++
				if cover == nil {
					cover, err = os.ReadFile(s.originalPath + "/" + entity.Name() + "/" + volumeEntity.Name() + "/0.jpeg")
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
