package main

import (
	"fmt"
	"os"

	"github.com/coderianx/cloudra/internal/client"
	"github.com/coderianx/cloudra/internal/config"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("  cloudra server <url>")
		fmt.Println("  cloudra upload <file>")
		fmt.Println("  cloudra download <file>")
		fmt.Println("  cloudra list")
		return
	}

	cmd := os.Args[1]

	switch cmd {

	// server set
	case "server":
		if len(os.Args) < 3 {
			fmt.Println("missing server url")
			return
		}
		config.SetServerURL(os.Args[2])
		fmt.Println("server set:", os.Args[2])

	// upload
	case "upload":
		if len(os.Args) < 3 {
			fmt.Println("missing file")
			return
		}
		err := client.Upload(os.Args[2])
		if err != nil {
			fmt.Println("upload error:", err)
		}

	// download
	case "download":
		if len(os.Args) < 3 {
			fmt.Println("missing file")
			return
		}
		err := client.Download(os.Args[2])
		if err != nil {
			fmt.Println("download error:", err)
		}

	// list
	case "list":
		err := client.List()
		if err != nil {
			fmt.Println("list error:", err)
		}

	default:
		fmt.Println("unknown command")
	}
}
