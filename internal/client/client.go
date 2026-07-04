package client

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coderianx/cloudra/internal/config"
)

func Upload(path string, recursive bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() && !recursive {
		return fmt.Errorf("%q is a directory, use -r to upload recursively", path)
	}

	if recursive && !info.IsDir() {
		return fmt.Errorf("-r expects a directory, got a file")
	}

	if recursive {
		return uploadDir(path)
	}

	return uploadFile(path)
}

func uploadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return err
	}

	io.Copy(part, file)
	writer.Close()

	_, err = http.Post(config.ServerURL+"/upload",
		writer.FormDataContentType(),
		&body,
	)

	return err
}

func uploadDir(dir string) error {
	tmpFile, err := os.CreateTemp("", "cloudra-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	zw := zip.NewWriter(tmpFile)

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if rel == "." {
				return nil
			}
			_, err := zw.Create(rel + "/")
			return err
		}

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		zw.Close()
		return err
	}

	zw.Close()
	tmpFile.Close()

	zipName := filepath.Base(dir) + ".zip"

	zipFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	defer zipFile.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", zipName)
	if err != nil {
		return err
	}

	io.Copy(part, zipFile)
	writer.Close()

	resp, err := http.Post(config.ServerURL+"/upload",
		writer.FormDataContentType(),
		&body,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	return nil
}

func Download(name string, zipMode bool) error {
	resp, err := http.Get(config.ServerURL + "/download?name=" + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server: %s", string(body))
	}

	if zipMode {
		out, err := os.Create(name + ".zip")
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		return err
	}

	tmpFile, err := os.CreateTemp("", "cloudra-*.zip")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	_, err = io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	if err != nil {
		return err
	}

	reader, zipErr := zip.OpenReader(tmpName)
	if zipErr != nil {
		out, err := os.Create(name)
		if err != nil {
			return err
		}
		defer out.Close()

		src, err := os.Open(tmpName)
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(out, src)
		return err
	}
	defer reader.Close()

	dirName := name
	if dirName == "" {
		dirName = "."
	}

	for _, f := range reader.File {
		fpath := filepath.Join(dirName, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		rc, err := f.Open()
		if err != nil {
			continue
		}

		out, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			continue
		}

		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}

	fmt.Println("extracted:", dirName)
	return nil
}

func List() error {
	resp, err := http.Get(config.ServerURL + "/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}
