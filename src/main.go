package main

import (
	"fmt"
	"time"

	"github.com/Mine4x/OpenLinkHub-system-tray/src/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	tray := systray.New("OpenLinkHub", "OpenLinkHub-sytem-tray", icon.Data)

	var testBool = false

	statusItem := tray.AddMenuItem("Status: True", "Current testBool value", nil)
	statusItem.SetEnabled(false)

	tray.AddSeparator()

	tray.AddMenuItem("Change bool", "Change testBool", func() {
		fmt.Println("Changing testBool!")
		testBool = !testBool
		if testBool == true {
			statusItem.SetTitle("Status: True")
		} else {
			statusItem.SetTitle("Status: False")
		}
	})

	testMenu := tray.AddMenuItem("Test Menu", "Test", nil)

	testMenu.AddSubMenuItem("Change bool", "Change testBool", nil)

	tray.OnReady(func() {
		fmt.Println("System tray is ready!")

		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				fmt.Printf("Test")
			}
		}()
	})

	tray.OnExit(func() {
		fmt.Printf("Exiting")
	})

	tray.Run()
}
