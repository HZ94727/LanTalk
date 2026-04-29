package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	winoptions "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist/index.html frontend/dist/assets/*
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.tray.Register()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "LanTalk",
		Width:     1080,
		Height:    740,
		MinWidth:  980,
		MinHeight: 720,
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "LanTalk.SingleInstance",
			OnSecondInstanceLaunch: app.handleSecondInstanceLaunch,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 247, B: 249, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnBeforeClose:    app.beforeClose,
		Windows: &winoptions.Options{
			Theme: winoptions.Light,
			CustomTheme: &winoptions.ThemeSettings{
				LightModeTitleBar:          winoptions.RGB(245, 247, 249),
				LightModeTitleBarInactive:  winoptions.RGB(239, 242, 246),
				LightModeTitleText:         winoptions.RGB(24, 30, 38),
				LightModeTitleTextInactive: winoptions.RGB(104, 116, 130),
				LightModeBorder:            winoptions.RGB(214, 220, 228),
				LightModeBorderInactive:    winoptions.RGB(226, 231, 237),
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
