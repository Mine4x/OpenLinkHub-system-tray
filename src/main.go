package main

import (
	"fmt"
	"time"

	"github.com/Mine4x/OpenLinkHub-system-tray/src/battery"
	"github.com/Mine4x/OpenLinkHub-system-tray/src/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	tray := systray.New("OpenLinkHub", "OpenLinkHub-system-tray", icon.Data)

	batteryItem := tray.AddMenuItem("Battery: Loading...", "Device battery information", nil)
	batteryItem.SetEnabled(false)

	batteryIcons, err := battery.GetIcons()
	if err != nil {
		fmt.Printf("Error getting battery icons: %w", err)
	}

	tray.AddSeparator()

	tray.AddMenuItem("Refresh", "Refresh battery information", func() {
		updateBatteryInfo(tray, batteryItem, batteryIcons) // FIXME: Delayed
	})

	tray.AddSeparator()

	tray.AddMenuItem("Quit", "Exit the application", func() {
		tray.Quit() // FIXME: Delayed
	})

	tray.OnReady(func() {
		fmt.Println("System tray is ready!")

		updateBatteryInfo(tray, batteryItem, batteryIcons)

		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				updateBatteryInfo(tray, batteryItem, batteryIcons)
			}
		}()
	})

	tray.OnExit(func() {
		fmt.Println("Cleaning up...")
	})

	fmt.Println("Starting OpenLinkHub system tray...")
	tray.Run()
}

func updateBatteryInfo(tray *systray.Tray, batteryItem *systray.MenuItem, _batteryIcons *battery.BatteryIcons) {
	stats, err := battery.GetBatteryStats()
	if err != nil {
		fmt.Printf("Error fetching battery stats: %v\n", err)
		batteryItem.SetTitle("Battery: Error")
		tray.SetTooltip("OpenLinkHub - Error fetching battery")
		return
	}

	if len(stats.Data) == 0 {
		batteryItem.SetTitle("Battery: No devices")
		tray.SetTooltip("OpenLinkHub - No devices found")
		fmt.Println("No devices found")
		return
	}

	var lowestDevice *battery.BatteryDevice
	var lowestSerial string
	lowestLevel := 101

	for serial, device := range stats.Data {
		if device.Level < lowestLevel {
			lowestLevel = device.Level
			deviceCopy := device
			lowestDevice = &deviceCopy
			lowestSerial = serial
		}
	}

	if lowestDevice != nil {
		title := fmt.Sprintf("🔋 %s: %d%%", lowestDevice.Device, lowestDevice.Level)
		batteryItem.SetTitle(title)

		tooltip := "OpenLinkHub Devices:\n"
		for serial, device := range stats.Data {
			tooltip += fmt.Sprintf("%s: %d%%\n", device.Device, device.Level)
			if serial == lowestSerial {
				tooltip += "(lowest)"
			}
		}
		tray.SetTooltip(tooltip)

		fmt.Printf("Battery updated: %s at %d%%\n", lowestDevice.Device, lowestDevice.Level)

		if lowestDevice.Level <= 20 {
			fmt.Printf("Low battery warning: %s at %d%%\n", lowestDevice.Device, lowestDevice.Level)
		}
	}
}
