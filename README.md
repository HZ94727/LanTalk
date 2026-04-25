# LanTalk

LanTalk is a local-network desktop messenger built with `Wails + Go + Vite`.

## What It Does

- Discovers peers on the same LAN with UDP broadcast
- Uses TCP for direct text messaging
- Stores your nickname and chat history locally in a JSON file
- Updates the UI in real time through Wails events

## Project Structure

- `app.go`: Wails bindings exposed to the frontend
- `internal/chat`: LAN discovery, direct messaging, and local persistence
- `frontend/src`: the chat UI

## Development

Install frontend dependencies once:

```bash
cd frontend
npm install
```

Run the desktop app in development mode:

```bash
wails dev
```

If `wails` is not in your `PATH`, on this machine it was installed with:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Build

```bash
wails build
```

The packaged executable is written to `build/bin/LanTalk.exe`.

## Local Data

LanTalk stores state in your user config directory:

- Windows: `%AppData%\\LanTalk\\state.json`

The app also exposes the exact path in the sidebar so it is easy to inspect.

## Current Scope

This first version focuses on:

- nickname management
- automatic peer discovery
- one-to-one text chat
- local conversation persistence

## Known Limitation

Discovery currently listens on a fixed UDP port, so testing two LanTalk instances on the same Windows machine is not the ideal path. The app is meant to be opened on two devices within the same LAN for realistic testing.
