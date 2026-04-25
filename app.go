package main

import (
	"context"

	"LanTalk/internal/chat"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	service *chat.Service
}

func NewApp() *App {
	app := &App{}

	service, err := chat.NewService("LanTalk", chat.Callbacks{
		OnProfileChanged: func(profile chat.Profile) {
			if app.ctx != nil {
				runtime.EventsEmit(app.ctx, "profile:updated", profile)
			}
		},
		OnSettingsChanged: func(settings chat.Settings) {
			if app.ctx != nil {
				runtime.EventsEmit(app.ctx, "settings:updated", settings)
			}
		},
		OnDataPathChanged: func(path string) {
			if app.ctx != nil {
				runtime.EventsEmit(app.ctx, "data:pathUpdated", path)
			}
		},
		OnPeersChanged: func(peers []chat.Peer) {
			if app.ctx != nil {
				runtime.EventsEmit(app.ctx, "peers:updated", peers)
			}
		},
		OnMessage: func(peerID string, messages []chat.ChatMessage) {
			if app.ctx != nil {
				runtime.EventsEmit(app.ctx, "conversation:updated", map[string]interface{}{
					"peerId":   peerID,
					"messages": messages,
				})
			}
		},
	})
	if err != nil {
		panic(err)
	}

	app.service = service
	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.service.Start(ctx); err != nil {
		panic(err)
	}
}

func (a *App) shutdown(context.Context) {
	a.service.Stop()
}

func (a *App) Bootstrap() chat.Snapshot {
	return a.service.Snapshot()
}

func (a *App) UpdateDisplayName(name string) (chat.Profile, error) {
	return a.service.UpdateDisplayName(name)
}

func (a *App) SendChatMessage(peerID, text string) error {
	return a.service.SendMessage(peerID, text)
}

func (a *App) SendImageMessage(peerID, dataURL, fileName string) error {
	return a.service.SendImageMessage(peerID, dataURL, fileName)
}

func (a *App) UpdateLanguage(language string) (chat.Settings, error) {
	return a.service.UpdateLanguage(language)
}

func (a *App) UpdateTheme(theme string) (chat.Settings, error) {
	return a.service.UpdateTheme(theme)
}

func (a *App) DataPath() string {
	return a.service.DataPath()
}

func (a *App) ChooseDataDirectory() (string, error) {
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose LanTalk storage directory",
	})
	if err != nil {
		return "", err
	}
	if selectedDir == "" {
		return a.service.DataPath(), nil
	}

	return a.service.MoveDataFile(selectedDir)
}

func (a *App) EnsureDebugPeer() chat.Peer {
	return a.service.EnsureDebugPeer()
}

func (a *App) AddManualPeer(name, address string, port int) (chat.Peer, error) {
	return a.service.AddManualPeer(name, address, port)
}
