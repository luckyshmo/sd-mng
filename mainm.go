package main

import (
	"fmt"

	"kek.com/cmd/formats"
	"kek.com/cmd/formats/download"
	md "kek.com/mangadex"
)

// const id = "10a4985d-0713-462e-a9d6-767bf91e4fd7"

func getMangaInfo(id string) (*md.Manga, error) {
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

	covers, err := getCovers(manga)
	if err != nil {
		return nil, fmt.Errorf("get covers: %w", err)
	}
	*manga = manga.WithCovers(covers)

	// for _, volume := range manga.Sorted() {
	// 	if err := handleVolume(*manga, volume); err != nil {
	// 		return nil, fmt.Errorf("volume %v: %w", volume.Info.Identifier, err)
	// 	}
	// }

	return manga, nil
}
