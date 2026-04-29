package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync/atomic"

	"LanTalk/internal/chat"
	"LanTalk/internal/clipimg"
	"LanTalk/internal/mediautil"
	"github.com/wailsapp/wails/v2/pkg/options"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	service    *chat.Service
	tray       *trayManager
	allowClose atomic.Bool
}

func NewApp() *App {
	app := &App{}

	service, err := chat.NewService("LanTalk", chat.Callbacks{
		OnProfileChanged: func(profile chat.Profile) {
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "profile:updated", profile)
			}
		},
		OnSettingsChanged: func(settings chat.Settings) {
			if app.tray != nil {
				app.tray.SetLanguage(settings.Language)
			}
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "settings:updated", settings)
			}
		},
		OnDataPathChanged: func(path string) {
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "data:pathUpdated", path)
			}
		},
		OnPeersChanged: func(peers []chat.Peer) {
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "peers:updated", peers)
			}
		},
		OnMessage: func(peerID string, messages []chat.ChatMessage) {
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "conversation:updated", map[string]interface{}{
					"peerId":   peerID,
					"messages": messages,
				})
			}
		},
		OnTransferProgress: func(progress chat.TransferProgress) {
			if app.ctx != nil {
				wruntime.EventsEmit(app.ctx, "transfer:progress", progress)
			}
		},
	})
	if err != nil {
		panic(err)
	}

	app.service = service
	app.tray = newTrayManager(app, service.Snapshot().Settings.Language)
	return app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.setMainWindowVisible(true)
	if err := a.service.Start(ctx); err != nil {
		panic(err)
	}
}

func (a *App) shutdown(context.Context) {
	a.allowClose.Store(true)
	if a.tray != nil {
		a.tray.Shutdown()
	}
	a.service.Stop()
}

func (a *App) beforeClose(context.Context) bool {
	if a.allowClose.Load() {
		return false
	}
	if a.tray == nil || !a.tray.IsReady() {
		return false
	}
	a.hideMainWindow()
	return true
}

func (a *App) Bootstrap() chat.Snapshot {
	return a.service.Snapshot()
}

func (a *App) LoadConversationPage(peerID string, beforeTimestamp int64, beforeID string, limit int) (chat.ConversationPage, error) {
	return a.service.LoadConversation(peerID, beforeTimestamp, beforeID, limit)
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

func (a *App) SendFileMessage(peerID, dataURL, fileName string) error {
	return a.service.SendFileMessage(peerID, dataURL, fileName)
}

func (a *App) SendLocalFileMessage(peerID, filePath string) error {
	return a.service.SendLocalFileMessage(peerID, filePath)
}

func (a *App) RetryMessage(peerID, messageID string) error {
	return a.service.RetryMessage(peerID, messageID)
}

func (a *App) DeleteMessage(peerID, messageID string) error {
	return a.service.DeleteMessage(peerID, messageID)
}

func (a *App) ClearConversation(peerID string) (int, error) {
	return a.service.ClearConversation(peerID)
}

func (a *App) LoadStorageStats() (chat.StorageStats, error) {
	return a.service.LoadStorageStats()
}

func (a *App) CleanupUnusedMedia() (chat.MediaCleanupResult, error) {
	return a.service.CleanupUnusedMedia()
}

func (a *App) CopyImageToClipboard(dataURL string) error {
	return clipimg.WriteImageDataURL(dataURL)
}

func (a *App) LoadImageSource(value, mediaType string) (string, error) {
	return a.service.ResolveImageDataURL(value, mediaType)
}

func (a *App) LoadFileSource(value, mediaType string) (string, error) {
	return a.service.ResolveFileDataURL(value, mediaType)
}

func (a *App) CheckFileMessageAvailable(value string) bool {
	return a.service.IsStoredFileAvailable(value)
}

func (a *App) SaveImageMessage(value, suggestedName string) (string, error) {
	raw, mediaType, err := a.service.ResolveImagePayload(value, "")
	if err != nil {
		return "", err
	}

	filename := mediautil.NormalizeImageFileName(suggestedName, mediaType)
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext == "" {
		ext = strings.TrimPrefix(mediautil.ExtensionForMediaType(mediaType), ".")
	}

	targetPath, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:           "Save image",
		DefaultFilename: filename,
		Filters: []wruntime.FileFilter{
			{
				DisplayName: "Image Files",
				Pattern:     fmt.Sprintf("*.%s", ext),
			},
			{
				DisplayName: "All Files",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if targetPath == "" {
		return "", nil
	}

	if err := os.WriteFile(targetPath, raw, 0o644); err != nil {
		return "", err
	}

	return targetPath, nil
}

func (a *App) SaveFileMessage(value, mediaType, suggestedName string) (string, error) {
	raw, resolvedMediaType, err := a.service.ResolveFilePayload(value, mediaType)
	if err != nil {
		return "", err
	}

	filename := mediautil.NormalizeFileName(suggestedName, resolvedMediaType)
	targetPath, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:           "Save file",
		DefaultFilename: filename,
		Filters: []wruntime.FileFilter{
			{
				DisplayName: "All Files",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if targetPath == "" {
		return "", nil
	}

	if err := os.WriteFile(targetPath, raw, 0o644); err != nil {
		return "", err
	}

	return targetPath, nil
}

func (a *App) OpenFileMessage(value string) error {
	targetPath, err := a.service.ResolveStoredFilePath(value)
	if err != nil {
		return err
	}
	if _, err := os.Stat(targetPath); err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", targetPath)
	case "darwin":
		cmd = exec.Command("open", targetPath)
	default:
		cmd = exec.Command("xdg-open", targetPath)
	}
	return cmd.Start()
}

func revealPathInFileManager(targetPath string) error {
	selectTarget := true
	if _, err := os.Stat(targetPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			selectTarget = false
			targetPath = filepath.Dir(targetPath)
		} else {
			return err
		}
	}

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		if selectTarget {
			cmd = exec.Command("explorer", "/select,"+targetPath)
		} else {
			cmd = exec.Command("explorer", targetPath)
		}
	case "darwin":
		if selectTarget {
			cmd = exec.Command("open", "-R", targetPath)
		} else {
			cmd = exec.Command("open", targetPath)
		}
	default:
		if selectTarget {
			targetPath = filepath.Dir(targetPath)
		}
		cmd = exec.Command("xdg-open", targetPath)
	}
	return cmd.Start()
}

func (a *App) RevealFileMessage(value string) error {
	targetPath, err := a.service.ResolveStoredFilePath(value)
	if err != nil {
		return err
	}

	return revealPathInFileManager(targetPath)
}

func (a *App) RevealImageMessage(value string) error {
	targetPath, err := a.service.ResolveStoredImagePath(value)
	if err != nil {
		return err
	}

	return revealPathInFileManager(targetPath)
}

func (a *App) SetUnreadCount(count int) {
	if a.tray != nil {
		a.tray.SetUnreadCount(count)
	}
}

func (a *App) setMainWindowVisible(visible bool) {
	if a.tray != nil {
		a.tray.SetWindowVisible(visible)
	}
	if a.ctx != nil {
		wruntime.EventsEmit(a.ctx, "window:visibilityChanged", map[string]bool{
			"visible": visible,
		})
	}
}

func (a *App) hideMainWindow() {
	if a.ctx == nil {
		return
	}
	wruntime.WindowHide(a.ctx)
	a.setMainWindowVisible(false)
}

func (a *App) showMainWindow() {
	if a.ctx == nil {
		return
	}
	wruntime.WindowUnminimise(a.ctx)
	wruntime.WindowShow(a.ctx)
	a.setMainWindowVisible(true)
}

func (a *App) requestQuit() {
	a.allowClose.Store(true)
	if a.ctx != nil {
		wruntime.Quit(a.ctx)
	}
}

func (a *App) handleSecondInstanceLaunch(options.SecondInstanceData) {
	a.showMainWindow()
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
	selectedDir, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
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
