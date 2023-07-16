package main

import (
	"fmt"
	"log"

	"kek.com/cmd/formats"
	"kek.com/cmd/formats/download"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const id = "10a4985d-0713-462e-a9d6-767bf91e4fd7"

func run() error {
	manga, err := download.MangadexSkeleton(id)
	if err != nil {
		return fmt.Errorf("skeleton: %w", err)
	}

	chapters, err := getChapters(*manga, id)
	if err != nil {
		return fmt.Errorf("chapters: %w", err)
	}
	*manga = manga.WithChapters(chapters)

	formats.PrintSummary(manga)

	covers, err := getCovers(manga)
	if err != nil {
		return fmt.Errorf("covers: %w", err)
	}
	*manga = manga.WithCovers(covers)

	for _, volume := range manga.Sorted() {
		if err := handleVolume(*manga, volume); err != nil {
			return fmt.Errorf("volume %v: %w", volume.Info.Identifier, err)
		}
	}

	return nil
}
