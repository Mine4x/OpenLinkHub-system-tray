package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("My App")
	systray.SetTooltip("My Go Application")

	mShow := systray.AddMenuItem("Show", "Show the window")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}
