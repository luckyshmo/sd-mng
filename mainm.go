package main

import (
	"fmt"

	"github.com/google/uuid"
	"kek.com/cmd/formats"
	"kek.com/cmd/formats/download"
	md "kek.com/mangadex"
)

// const id = "10a4985d-0713-462e-a9d6-767bf91e4fd7"

type MangaPreviewInfo struct {
	Info    md.MangaInfo
	Volumes []md.VolumeSorted
}

func getMangaInfo(id string) (*MangaPreviewInfo, error) {
	manga, err := download.MangadexSkeleton(id)
	if err != nil {
		return nil, fmt.Errorf("download skeleton: %w", err)
	}

	chapters, err := getChapters(manga, id)
	if err != nil {
		return nil, fmt.Errorf("get chapters: %w", err)
	}
	manga = manga.WithChapters(chapters)

	formats.PrintSummary(manga)

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
	return &MangaPreviewInfo{manga.Info, sortedVolumes}, nil

	// covers, err := getCovers(manga)
	// if err != nil {
	// 	return nil, fmt.Errorf("get covers: %w", err)
	// }
	// *manga = manga.WithCovers(covers)

	// for _, volume := range manga.Sorted() {
	// 	if err := handleVolume(*manga, volume); err != nil {
	// 		return nil, fmt.Errorf("volume %v: %w", volume.Info.Identifier, err)
	// 	}
	// }
}
