package main

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/getlantern/systray"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

type trayStrings struct {
	appName      string
	showWindow   string
	hideWindow   string
	quit         string
	noUnread     string
	unreadFormat string
	tooltipBase  string
	tooltipCount string
}

func trayCopy(language string) trayStrings {
	if language == "en-US" {
		return trayStrings{
			appName:      "LanTalk",
			showWindow:   "Show window",
			hideWindow:   "Hide to tray",
			quit:         "Quit",
			noUnread:     "No unread messages",
			unreadFormat: "Unread messages: %d",
			tooltipBase:  "LanTalk",
			tooltipCount: "LanTalk - %d unread",
		}
	}

	return trayStrings{
		appName:      "LanTalk",
		showWindow:   "显示主界面",
		hideWindow:   "隐藏到托盘",
		quit:         "退出",
		noUnread:     "当前没有未读消息",
		unreadFormat: "未读消息：%d",
		tooltipBase:  "LanTalk",
		tooltipCount: "LanTalk - %d 条未读消息",
	}
}

type trayManager struct {
	app *App

	mu            sync.RWMutex
	ready         bool
	windowVisible bool
	unreadCount   int
	language      string

	titleItem  *systray.MenuItem
	unreadItem *systray.MenuItem
	toggleItem *systray.MenuItem
	quitItem   *systray.MenuItem
}

func newTrayManager(app *App, language string) *trayManager {
	return &trayManager{
		app:           app,
		windowVisible: true,
		language:      normalizeTrayLanguage(language),
	}
}

func normalizeTrayLanguage(language string) string {
	if language == "en-US" {
		return "en-US"
	}
	return "zh-CN"
}

func (t *trayManager) Register() {
	systray.Register(t.onReady, t.onExit)
}

func (t *trayManager) Shutdown() {
	systray.Quit()
}

func (t *trayManager) IsReady() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ready
}

func (t *trayManager) SetUnreadCount(count int) {
	if count < 0 {
		count = 0
	}
	t.mu.Lock()
	t.unreadCount = count
	t.refreshLocked()
	t.mu.Unlock()
}

func (t *trayManager) SetWindowVisible(visible bool) {
	t.mu.Lock()
	t.windowVisible = visible
	t.refreshLocked()
	t.mu.Unlock()
}

func (t *trayManager) SetLanguage(language string) {
	t.mu.Lock()
	t.language = normalizeTrayLanguage(language)
	t.refreshLocked()
	t.mu.Unlock()
}

func (t *trayManager) onReady() {
	if len(trayIcon) > 0 {
		systray.SetTemplateIcon(trayIcon, trayIcon)
	}

	titleItem := systray.AddMenuItem("LanTalk", "LanTalk")
	titleItem.Disable()
	unreadItem := systray.AddMenuItem("", "")
	unreadItem.Disable()
	systray.AddSeparator()
	toggleItem := systray.AddMenuItem("", "")
	quitItem := systray.AddMenuItem("", "")

	t.mu.Lock()
	t.titleItem = titleItem
	t.unreadItem = unreadItem
	t.toggleItem = toggleItem
	t.quitItem = quitItem
	t.ready = true
	t.refreshLocked()
	t.mu.Unlock()

	go t.eventLoop(toggleItem, quitItem)
}

func (t *trayManager) onExit() {}

func (t *trayManager) eventLoop(toggleItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-toggleItem.ClickedCh:
			t.mu.RLock()
			visible := t.windowVisible
			t.mu.RUnlock()
			if visible {
				t.app.hideMainWindow()
			} else {
				t.app.showMainWindow()
			}
		case <-quitItem.ClickedCh:
			t.app.requestQuit()
			return
		}
	}
}

func (t *trayManager) refreshLocked() {
	copy := trayCopy(t.language)

	systray.SetTooltip(t.tooltipLocked(copy))

	if !t.ready {
		return
	}

	if t.titleItem != nil {
		t.titleItem.SetTitle(copy.appName)
	}
	if t.unreadItem != nil {
		t.unreadItem.SetTitle(t.unreadLabelLocked(copy))
	}
	if t.toggleItem != nil {
		if t.windowVisible {
			t.toggleItem.SetTitle(copy.hideWindow)
		} else {
			t.toggleItem.SetTitle(copy.showWindow)
		}
	}
	if t.quitItem != nil {
		t.quitItem.SetTitle(copy.quit)
	}
}

func (t *trayManager) unreadLabelLocked(copy trayStrings) string {
	if t.unreadCount <= 0 {
		return copy.noUnread
	}
	return fmt.Sprintf(copy.unreadFormat, t.unreadCount)
}

func (t *trayManager) tooltipLocked(copy trayStrings) string {
	if t.unreadCount <= 0 {
		return copy.tooltipBase
	}
	return fmt.Sprintf(copy.tooltipCount, t.unreadCount)
}
