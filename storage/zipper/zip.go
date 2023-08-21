package zipper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func ZipFolder(rootFolder string) error {
	Volumes, err := os.ReadDir(rootFolder)
	if err != nil {
		return fmt.Errorf("error reading root folder: %w", err)
	}

	for _, volume := range Volumes {
		if volume.IsDir() {
			processVolume(volume.Name(), rootFolder+"/"+volume.Name(), "./dist")
		}
	}

	return nil
}

type Image struct {
	relativePath string
	img          []byte
}

func processChapter(name, volumeName, rootFolder string) ([]Image, error) {
	fmt.Println("want to open: ", rootFolder)
	images, err := os.ReadDir(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	var res = make([]Image, len(images))
	for i, img := range images {
		if img.IsDir() {
			log.Fatal("ERROR: should not be any dirs")
		}
		image, err := os.Open(path.Join(rootFolder, img.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer image.Close()
		data, err := io.ReadAll(image)
		if err != nil {
			log.Fatal(err)
		}
		res[i] = Image{
			volumeName + "/" + name + "/" + name + "_" + img.Name(),
			data,
		}

	}

	return res, nil
}

func processVolume(volumeName, rootFolder, distPath string) error {
	fmt.Println("Processing ", volumeName)

	err := os.MkdirAll(distPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	chapters, err := os.ReadDir(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	var imgPaths []Image
	for _, chapter := range chapters {
		if chapter.IsDir() {
			cip, err := processChapter(chapter.Name(), volumeName, rootFolder+"/"+chapter.Name())
			if err != nil {
				log.Fatal("kek: ", err)
			}
			imgPaths = append(imgPaths, cip...)
			continue
		}
		cover, err := os.Open(path.Join(rootFolder, chapter.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer cover.Close()
		data, err := io.ReadAll(cover)
		if err != nil {
			log.Fatal(err)
		}
		// add cover to copy list
		imgPaths = append(imgPaths, Image{
			chapter.Name(),
			data,
		})
	}

	// Create cbz for cbz
	cbzFile := fmt.Sprintf("%s.cbz", path.Join(distPath, volumeName))
	zipFile, err := os.Create(cbzFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, img := range imgPaths {
		fmt.Println("img path: ", img.relativePath)
		writer, err := zipWriter.Create(img.relativePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, bytes.NewReader(img.img))
		if err != nil {
			return err
		}
	}

	return nil
}
