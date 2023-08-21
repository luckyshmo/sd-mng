package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"kek.com/cmd/formats/download"
	md "kek.com/mangadex"
	"kek.com/storage"

	"kek.com/cache"
	"kek.com/upscaler"
)

type MangaPreviewInfo struct {
	Info    md.MangaInfo
	Volumes []md.VolumeSorted
}

type MangaUC struct {
	upscaler upscaler.Upscaler
	storage  storage.Storage
}

func NewMangaUC(ups upscaler.Upscaler, storage storage.Storage) *MangaUC {
	return &MangaUC{ups, storage}
}

func getMangaInfo(id string) (*MangaPreviewInfo, error) {
	m, ok := cache.AppCache.Get(id)
	if ok {
		fmt.Println("hit")
		return m.(*MangaPreviewInfo), nil
	}

	manga, err := download.MangadexSkeleton(id)
	if err != nil {
		return nil, fmt.Errorf("download skeleton: %w", err)
	}

	chapters, err := getChapters(manga, id)
	if err != nil {
		return nil, fmt.Errorf("get chapters: %w", err)
	}
	manga = manga.WithChapters(chapters)

	coverPaths, err := getCoverPaths(manga)
	if err != nil {
		return nil, fmt.Errorf("get covers: %w", err)
	}

	fmt.Println(coverPaths)
	manga = manga.WithCoverPaths(coverPaths)

	volumes := manga.Sorted()
	sortedVolumes := make([]md.VolumeSorted, len(volumes))
	for i, v := range volumes {
		sortedVolumes[i] = md.VolumeSorted{
			UID:       uuid.NewString(),
			Info:      v.Info,
			Chapters:  v.Sorted(),
			CoverPath: v.CoverPath,
		}
	}

	mInfo := &MangaPreviewInfo{manga.Info, sortedVolumes}
	cache.AppCache.Set(id, mInfo)

	return mInfo, nil
}

func downloadManga(id string, idsToAdd map[string]map[string]interface{}) error {
	m, ok := cache.AppCache.Get(id)
	if !ok {
		return fmt.Errorf("session not found") //!
	}

	previewInfo := m.(*MangaPreviewInfo)

	var manga md.Manga
	manga.Volumes = make(map[md.Identifier]md.Volume)
	manga.Info = previewInfo.Info
	for _, volumePreview := range previewInfo.Volumes {
		idsToAdd, ok := idsToAdd[volumePreview.UID]
		if !ok {
			continue
		}
		var vol md.Volume
		vol.Chapters = make(map[md.Identifier]md.Chapter)
		vol.Info = volumePreview.Info
		for _, chapterPreview := range volumePreview.Chapters {
			if _, ok := idsToAdd[chapterPreview.UID]; ok {
				vol.Chapters[chapterPreview.Info.Identifier] = chapterPreview
			}
		}
		if len(vol.Chapters) > 0 {
			manga.Volumes[vol.Info.Identifier] = vol
		}
	}

	covers, err := getCovers(&manga)
	if err != nil {
		return fmt.Errorf("get covers: %w", err)
	}
	manga = manga.WithCovers(covers) //!

	for _, volume := range manga.Sorted() {
		if err := handleVolume(manga, volume); err != nil {
			return fmt.Errorf("volume %v: %w", volume.Info.Identifier, err)
		}
	}
	fmt.Println("done")
	return nil
}

func (md *MangaUC) Upscale(title string) error {
	reader := md.storage.NewReader(title)

	for img, ok := reader.Next(); ok; img, ok = reader.Next() {
		imgs, names, path := md.upscaler.Upscale(img)

		if len(imgs) == 0 {
			log.Println("no images")
			continue
		}

		if err := md.storage.Write(imgs, names, title, path); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	return nil
}

func (md *MangaUC) Zip(title string) error {
	return md.storage.Zip(title)
}
