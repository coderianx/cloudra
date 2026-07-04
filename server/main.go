package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const storageDir = "./storage"

func ensureStorage() {
	os.MkdirAll(storageDir, os.ModePerm)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", 405)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file error", 400)
		return
	}
	defer file.Close()

	name := header.Filename
	dstPath := filepath.Join(storageDir, name)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "cannot save file", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	fmt.Fprintln(w, "uploaded:", name)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name", 400)
		return
	}

	filePath := filepath.Join(storageDir, name)

	if _, err := os.Stat(filePath); err != nil {
		zipPath := filePath + ".zip"
		if _, err := os.Stat(zipPath); err == nil {
			http.ServeFile(w, r, zipPath)
			return
		}
		http.Error(w, "not found", 404)
		return
	}

	http.ServeFile(w, r, filePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(storageDir)
	if err != nil {
		http.Error(w, "cannot read storage", 500)
		return
	}

	for _, e := range entries {
		fmt.Fprintln(w, e.Name())
	}
}

func main() {
	ensureStorage()

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/list", listHandler)

	fmt.Println("Cloudra server running on :8080")
	http.ListenAndServe(":8080", nil)
}
