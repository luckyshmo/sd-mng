package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"
)

const (
	relativePath = "."
	loraPath     = "./models/Lora"
)

var id atomic.Uint64

// progressWriter is a custom io.Writer implementation that keeps track of the download progress.
type progressWriter struct {
	total      int
	downloaded int
	id         string
	name       string
	origin     string
	prev       int
	messenger  Messenger
}

func newProgressWriter(total int, id, name, origin string, messenger Messenger) *progressWriter {
	return &progressWriter{total, 0, id, name, origin, -1, messenger}
}

type Message struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}

func (pw *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	pw.downloaded += n
	percentage := (pw.downloaded * 100) / pw.total
	if pw.prev != percentage {
		msg, _ := json.Marshal(Message{pw.id, pw.name, percentage})
		pw.prev = percentage
		pw.messenger.Send(pw.origin, msg)
		// fmt.Println("send message from writer: ", string(msg))
		// fmt.Printf("\rDownloading... %d%%", (pw.downloaded*100)/pw.total)
	}
	return n, nil
}

type Downloader struct {
	messenger Messenger
}

func NewDownloader(m Messenger) *Downloader {
	return &Downloader{m}
}

type LoraInfo struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func (d *Downloader) LoraInfo() ([]LoraInfo, error) {
	// Read the folder
	files, err := os.ReadDir(loraPath)
	if err != nil {
		return nil, err
	}

	res := make([]LoraInfo, len(files))
	// Print the file names
	for i, file := range files {
		spl := strings.Split(file.Name(), ".")
		token := append([]string{"<lora:"}, spl[:len(spl)-1]...)
		token = append(token, ":1>")
		res[i] = LoraInfo{file.Name(), strings.Join(token, "")}
	}

	return res, nil
}

func (d *Downloader) DownloadFile(url, origin, folder string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	contentDisposition := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return fmt.Errorf("failed to parse Content-Disposition header: %s", err)
	}

	fileName := params["filename"]
	if fileName == "" {
		return fmt.Errorf("unable to extract file name from Content-Disposition header")
	}

	file, err := os.Create(path.Join(relativePath, folder, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the total file size
	fileSize := response.ContentLength

	// Create the progress writer
	idStr := fmt.Sprint(id.Load())
	progress := newProgressWriter(int(fileSize), idStr, fileName, origin, d.messenger) //&progressWriter{total: int(fileSize)}
	id.Add(1)

	// Create a multi-writer to write to both the output file and the progress writer
	writer := io.MultiWriter(file, progress)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\nfinished\n")

	return nil
}
