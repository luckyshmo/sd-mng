package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed fe-sdd/dist
var content embed.FS

var production = false

func main() {
	if production {
		go runFE()
	}
	ws := NewWebSockets()
	d := NewDownloader(ws)

	http.Handle("/upscale/info", corsHandler(magadexInfoHandler()))
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
