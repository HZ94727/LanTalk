# LanTalk

[简体中文](README.md) | English

LanTalk is a local-network desktop messenger built with `Wails + Go + Vite`. It is designed for computers on the same LAN and supports peer discovery, one-to-one chat, image transfer, file transfer, local history, and system tray integration.

## Features

- Discover LanTalk users on the same LAN through UDP broadcast
- Send text, image, and file messages over direct TCP connections
- Show message delivery states, retry failed messages, delete messages, forward messages, and clear conversations
- Preview, zoom, copy, save, and reveal image messages
- Save, open, reveal, and validate file messages
- Store contacts, settings, and chat history in SQLite
- Store image and file payloads on disk with hash-based deduplication
- Encrypt sensitive local fields with AES-GCM
- Switch between Simplified Chinese and English
- Switch between multiple UI themes
- Use system tray, hide-to-tray, unread count, and single-instance restore
- Use Echo Bot and manual peer connection for local testing

## Tech Stack

- Desktop framework: Wails v2
- Backend: Go
- Frontend build tool: Vite
- Frontend implementation: vanilla JavaScript, HTML, and CSS
- Local database: SQLite via `modernc.org/sqlite`
- LAN communication: UDP broadcast discovery and TCP peer-to-peer transfer
- System tray: `github.com/getlantern/systray`

## Project Structure

```text
.
├── app.go                    # Wails bindings
├── main.go                   # App entry point and window configuration
├── tray.go                   # System tray integration
├── internal
│   ├── chat                  # Chat, discovery, storage, encryption, and media handling
│   ├── clipimg               # Copy images to the clipboard
│   └── mediautil             # Media type, Data URL, and filename helpers
├── frontend
│   ├── src                   # Frontend UI and interactions
│   └── wailsjs               # Generated Wails bindings
├── build                     # Icons, manifests, and installer templates
├── PROJECT_OVERVIEW.md       # Detailed project progress document
└── wails.json                # Wails configuration
```

## Development

Install frontend dependencies:

```bash
cd frontend
npm install
```

Run in development mode:

```bash
wails dev
```

If Wails CLI is not installed:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Build

```bash
wails build
```

The Windows executable is usually generated at:

```text
build/bin/LanTalk.exe
```

## Local Data

LanTalk stores data in the user config directory by default:

```text
%AppData%\LanTalk
```

Main local data:

- `lantalk.db`: SQLite database
- `master.key`: local encryption key
- `media/`: image and file payloads

The settings panel shows the current storage directory and supports moving it to another location.

## Tests

```bash
go test ./...
```

## Current Limitations

- UDP discovery uses the fixed port `48555`, so running multiple real instances on the same machine is not the primary test path.
- LAN discovery depends on broadcast support. Firewalls, virtual adapters, or network isolation may affect discovery.
- File transfer is currently suitable for small attachments. Chunked large-file transfer, resume support, and trusted pairing are not implemented yet.

See [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md) for more details.
