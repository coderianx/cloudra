package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var ServerURL = "http://localhost:8080"

const configDirName = ".cloudra"
const configFileName = "config"

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName), nil
}

func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func Load() {
	path, err := ConfigFile()
	if err != nil {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var url string
	_, err = fmt.Sscanf(string(data), "server_url=%s", &url)
	if err != nil || url == "" {
		return
	}

	ServerURL = url
}

func SetServerURL(url string) error {
	ServerURL = url

	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := ConfigFile()
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte("server_url="+url+"\n"), 0644)
}
