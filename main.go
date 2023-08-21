package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"kek.com/storage"
	"kek.com/upscaler"
)

//go:embed fe-sdd/dist
var content embed.FS

var production = true

func main() {
	if production {
		go runFE()
	}
	ws := NewWebSockets()
	d := NewDownloader(ws)
	storage := storage.NewFSStorage()
	ups := upscaler.NewSDUpscaler(upscaler.UpscalerConfig{
		APIURL:       os.Getenv("UPSCALER_API_URL"),
		UpscalerType: os.Getenv("UPSCALER_TYPE"),
	})

	mangaUseCase := NewMangaUC(ups, storage)

	http.Handle("/manga/zip", corsHandler(zipHandler(mangaUseCase)))
	http.Handle("/manga/upscale", corsHandler(upscaleHandler(mangaUseCase)))
	http.Handle("/manga/origin/info", corsHandler(GetStoredMangaInfo(storage)))
	http.Handle("/storage/manga/", corsHandler(http.StripPrefix("/storage/manga/", http.FileServer(http.Dir(os.Getenv("MANGA_STORAGE_DIR"))))))
	http.Handle("/upscale/info", corsHandler(magadexInfoHandler()))
	http.Handle("/upscale/download", corsHandler(magadexDownloadHandler()))
	http.Handle("/", corsHandler(downloadHandler(d)))
	http.Handle("/progress", ws.ProgressHandler())
	http.Handle("/info/lora", corsHandler(loraInfoHandler(d)))

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func runFE() {
	buildFolder, err := fs.Sub(content, "fe-sdd/dist")
	if err != nil {
		log.Fatal(err)
	}

	httpFS := http.FS(buildFolder)

	fs := http.FileServer(httpFS)

	server := &http.Server{
		Addr:    ":3000",
		Handler: fs,
	}

	log.Println("Server started on http://localhost:3000")
	log.Fatal(server.ListenAndServe())
}

func GetStoredMangaInfo(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := s.GetStoredMangaInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func zipHandler(mangaUC *MangaUC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		if title == "" {
			http.Error(w, "Please provide a title to zip.", http.StatusBadRequest)
			return
		}

		err := mangaUC.Zip(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Zipped successfully.")
	}
}

func upscaleHandler(mangaUC *MangaUC) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		if title == "" {
			http.Error(w, "Please provide a title to upscale.", http.StatusBadRequest)
			return
		}

		err := mangaUC.Upscale(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Upscaled successfully.")
	}
}

func loraInfoHandler(d *Downloader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := d.LoraInfo()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error downloading file: %s", err), http.StatusInternalServerError)
			return
		}

		// Marshal the array into JSON
		jsonData, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the response
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func magadexInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Please provide a URL to download.", http.StatusBadRequest)
			return
		}

		mng, err := getMangaInfo(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(mng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the response
		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func magadexDownloadHandler() http.HandlerFunc {
	type DownloadRequest map[string]map[string]any

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Please provide a URL to download.", http.StatusBadRequest)
			return
		}

		var req DownloadRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(req) == 0 {
			http.Error(w, "Please provide a list of chapters to download.", http.StatusBadRequest)
			return
		}

		if downloadManga(id, req) != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Manga downloaded successfully.")
	}
}

func downloadHandler(d *Downloader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Please provide a URL to download.", http.StatusBadRequest)
			return
		}

		folder := r.URL.Query().Get("folder")
		if folder == "" {
			http.Error(w, "Please provide download folder path.", http.StatusBadRequest)
			return
		}

		origin := r.Header.Get("Origin")

		err := d.DownloadFile(url, origin, folder)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error downloading file: %s", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "File downloaded successfully.")
	}
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}
