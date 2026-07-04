package main

import (
	"fmt"
	"os"

	"github.com/coderianx/cloudra/internal/client"
	"github.com/coderianx/cloudra/internal/config"
)

func main() {
	config.Load()

	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  cloudra server <url>")
		fmt.Println("  cloudra upload [-r] <file|dir>")
		fmt.Println("  cloudra download [--zip] <file|dir>")
		fmt.Println("  cloudra list")
		return
	}

	cmd := os.Args[1]

	switch cmd {

	case "server":
		if len(os.Args) < 3 {
			fmt.Println("missing server url")
			return
		}
		err := config.SetServerURL(os.Args[2])
		if err != nil {
			fmt.Println("failed to save config:", err)
			return
		}
		fmt.Println("server set:", os.Args[2])

	case "upload":
		if len(os.Args) < 3 {
			fmt.Println("missing file or directory")
			return
		}

		recursive := false
		path := os.Args[2]

		if path == "-r" {
			recursive = true
			if len(os.Args) < 4 {
				fmt.Println("missing directory")
				return
			}
			path = os.Args[3]
		}

		err := client.Upload(path, recursive)
		if err != nil {
			fmt.Println("upload error:", err)
		}

	case "download":
		if len(os.Args) < 3 {
			fmt.Println("missing file or directory")
			return
		}

		zipMode := false
		path := os.Args[2]

		if path == "--zip" {
			zipMode = true
			if len(os.Args) < 4 {
				fmt.Println("missing file or directory")
				return
			}
			path = os.Args[3]
		}

		err := client.Download(path, zipMode)
		if err != nil {
			fmt.Println("download error:", err)
		}

	case "list":
		err := client.List()
		if err != nil {
			fmt.Println("list error:", err)
		}

	default:
		fmt.Println("unknown command")
	}
}
