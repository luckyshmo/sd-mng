package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"

	"golang.org/x/text/language"

	"kek.com/cmd/filter"
	"kek.com/cmd/formats/download"
	md "kek.com/mangadex"
)

func getChapters(manga *md.Manga, id string) (md.ChapterList, error) {
	chapters, err := download.MangadexChapters(id)
	if err != nil {
		return nil, fmt.Errorf("mangadex: %w", err)
	}

	chapters, err = sortFromFlags(chapters)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	return filter.RemoveDuplicates(chapters), nil
}

func getCovers(manga *md.Manga) (md.ImageList, error) {
	covers, err := download.MangadexCovers(manga)
	if err != nil {
		return nil, fmt.Errorf("mangadex: %w", err)
	}

	return covers, nil
}

func sortFromFlags(cl md.ChapterList) (md.ChapterList, error) {
	lang := language.Make("en")
	cl = filter.FilterByLanguage(cl, lang)
	// cl = filter.SortByMost(cl)
	cl = filter.SortByNewest(cl)

	return cl, nil
}

func handleVolume(skeleton md.Manga, volume md.Volume) error {
	pages, err := getPages(volume)
	if err != nil {
		return fmt.Errorf("pages: %w", err)
	}

	fmt.Println("Img info [0]: ")
	fmt.Println(pages[0].ImageIdentifier)
	fmt.Println(pages[0].VolumeIdentifier)
	fmt.Println(pages[0].ChapterIdentifier)

	fmt.Println("got ", len(pages), " pages")

	title := fmt.Sprintf("%v: %v",
		skeleton.Info.Title,
		volume.Info.Identifier.StringFilled(0, 0, false),
	)
	err = saveImage(volume.Cover, "./manga/"+title, 0, "jpeg")
	if err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	for _, p := range pages {
		// format, err := getImageFormat(p.Image)
		// if err != nil {
		// 	fmt.Printf("err determine format for image #%d in volume #%s", p.ImageIdentifier, volume.Info.Identifier.String())
		// 	continue
		// }
		format := "jpeg"

		err := saveImage(p.Image, "./manga/"+title+"/"+p.ChapterIdentifier.String(), p.ImageIdentifier+1, format)
		if err != nil {
			return fmt.Errorf("save image: %w", err)
		}
	}

	return nil
}

func getImageFormat(img image.Image) (string, error) {
	switch img.(type) {
	case *image.RGBA:
		return "png", nil
	case *image.YCbCr:
		return "jpeg", nil
	}

	return "", fmt.Errorf("unknown image format")
}

func saveImage(imageData image.Image, basePath string, filename int, format string) error {
	fileName := fmt.Sprintf("%d.%s", filename, format)

	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return fmt.Errorf("create volume path: %w", err)
	}

	imgPath := path.Join(basePath, fileName)
	// fmt.Println("trying to save: ", imgPath)

	file, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, imageData, nil)
	if err != nil {
		return err
	}

	// fmt.Printf("Image saved as: %s\n", fileName)
	return nil
}

func getPages(volume md.Volume) (md.ImageList, error) {
	mangadexPages, err := download.MangadexPages(volume.Sorted().FilterBy(func(ci md.ChapterInfo) bool {
		return ci.GroupNames.String() != "Filesystem"
	}), 0)
	if err != nil {
		return nil, fmt.Errorf("mangadex: %w", err)
	}

	return mangadexPages, nil
}
