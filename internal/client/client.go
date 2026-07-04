package client

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/coderianx/cloudra/internal/config"
)

func Upload(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
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

func Download(name string) error {
	resp, err := http.Get(config.ServerURL + "/download?name=" + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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
