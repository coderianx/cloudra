
<h1 align="center">
  <br>
  <img src=".github/logo.svg" alt="Cloudra" width="200">
  <br>
  Cloudra
  <br>
</h1>

<h4 align="center">Minimal, single-binary file transfer over HTTP — server & CLI client in pure Go.</h4>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.26.3-%2300ADD8?style=flat&logo=go">
  <img src="https://img.shields.io/badge/license-MIT-%23eba613?style=flat">
  <img src="https://img.shields.io/badge/deps-none-%2355dd55?style=flat">
  <img src="https://img.shields.io/badge/build-passing-%2333cc33?style=flat">
  <img src="https://img.shields.io/badge/PRs-welcome-%23ff69b4?style=flat">
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#usage">Usage</a> •
  <a href="#api">API</a> •
  <a href="#project-structure">Structure</a> •
  <a href="#building">Building</a> •
  <a href="#roadmap">Roadmap</a>
</p>

<br>

```
_________ .__                   .___              
\_   ___ \|  |   ____  __ __  __| _/___________   
/    \  \/|  |  /  _ \|  |  \/ __ |\_  __ \__  \  
\     \___|  |_(  <_> )  |  / /_/ | |  | \// __ \_
 \______  /____/\____/|____/\____ | |__|  (____  /
        \/                       \/            \/ 
   ⚡ Drop files. Anywhere. Instantly.
```

---

## Features

| | Feature | Detail |
|---|---------|--------|
| ⚡ | **Zero deps** | Pure Go standard library — no frameworks, no bloat. |
| 🪶 | **Single binary** | One binary for the server, one for the CLI. Drop & run. |
| 📤 | **Upload** | Push files via `cloudra upload <file>`. |
| 📥 | **Download** | Pull files via `cloudra download <file>`. |
| 📋 | **List** | See remote files with `cloudra list`. |
| 🔌 | **Pluggable server** | Point the CLI at any running Cloudra server. |
| 💾 | **Persistent config** | Server URL saved in `~/.cloudra/config`. |
| 📦 | **Directory upload** | Upload folders via `cloudra upload -r <dir>` — stored as `.zip` on server. |
| 📂 | **Directory download** | Download folders as `.zip` or extract locally with `cloudra download [--zip] <dir>`. |

---

## Quick Start

```bash
# Install the CLI directly
go install github.com/coderianx/cloudra/cmd/cloudra@latest

# Start a server from source
git clone https://github.com/coderianx/cloudra.git
cd cloudra
go run ./server/main.go &

# Upload a file
cloudra upload myfile.jpg

# Upload a directory (recursive — zips & extracts)
cloudra upload -r myfolder/

# List remote files
cloudra list

# Download a file
cloudra download myfile.jpg

# Download a directory (as a folder)
cloudra download myfolder/

# Download a directory (as .zip)
cloudra download --zip myfolder/
```

---

## Usage

### Server

```bash
# Start with default port (:8080)
go run ./server/main.go

# Files are stored in ./storage/ relative to the server binary
```

### CLI Client

```bash
# Set a custom server URL (saved to ~/.cloudra/config)
cloudra server http://192.168.1.100:8080

# Upload a file
cloudra upload photo.jpg

# Upload a directory (zipped & extracted on server)
cloudra upload -r documents/

# Download a file
cloudra download document.pdf

# Download a directory (extracted)
cloudra download documents/

# Download a directory (as .zip)
cloudra download --zip documents/

# List remote files
cloudra list
```

> Server URL is persisted in `~/.cloudra/config` — set it once, use it forever.

---

## API

### `POST /upload`
Upload a file via multipart form-data.

| Field | Type | Description |
|-------|------|-------------|
| `file` | file | The file to upload |

```
curl -F "file=@photo.jpg" http://localhost:8080/upload
```

### `GET /download?name=<filename>`
Download a previously uploaded file.

```
curl -O http://localhost:8080/download?name=photo.jpg
```

### `GET /list`
List all stored files.

```
curl http://localhost:8080/list
```

**Response:**
```json
["photo.jpg", "document.pdf"]
```

---

## Project Structure

```
cloudra/
├── cmd/
│   └── cloudra/
│       └── main.go             # CLI entry point
├── internal/
│   ├── client/
│   │   └── client.go           # HTTP client (upload / download / list)
│   └── config/
│       └── config.go           # Config file (~/.cloudra/config)
├── server/
│   └── main.go                 # HTTP server (3 routes, local storage)
├── go.mod                      # Go module definition
└── storage/                    # Created at runtime — file storage dir
```

---

## Building

```bash
# Build the server
go build -o bin/cloudra-server ./server/main.go

# Build the CLI client
go build -o bin/cloudra ./cmd/cloudra/main.go

# Build both
./scripts/build.sh
```

Binaries land in `./bin/`. No external dependencies required.

---

## Roadmap

- [ ] TLS support
- [ ] File-size limits & validation
- [ ] Authentication tokens
- [x] Directory upload (recursive — zip-based)
- [x] Directory download (as folder or .zip)
- [ ] Streaming progress bars
- [ ] Docker image
- [ ] GitHub Actions CI
- [ ] Unit & integration tests

---

## Contributing

Contributions are welcome! Open an issue or submit a PR.

1. Fork it
2. Create your feature branch (`git checkout -b feat/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing`)
5. Open a Pull Request

---

<p align="center">
  <sub>Built with ❤️ and the Go standard library.</sub>
  <br>
  <a href="https://github.com/coderianx/cloudra">github.com/coderianx/cloudra</a>
</p>
