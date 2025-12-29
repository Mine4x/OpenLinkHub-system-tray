package battery

import (
	"fmt"
	"sync"
	"time"

	"github.com/Mine4x/OpenLinkHub-system-tray/src/systray"
	"github.com/getlantern/systray/example/icon"
)

type traysEntry struct {
	serial      string
	batteryItem systray.MenuItem
	updated     bool
}

var mainTray *systray.Tray

func fetchDevices() (*BatteryResponse, error) {
	stats, err := GetBatteryStats()
	if err != nil {
		return nil, fmt.Errorf("Error fetching battery stats: %v", err)
	}
	if len(stats.Data) == 0 {
		return nil, fmt.Errorf("No devices")
	}
	return stats, nil
}

func handleDevice(serial string, device BatteryDevice, trays *[]traysEntry) error {
	battString := fmt.Sprintf("🔋 %s: %d%%", device.Device, device.Level)

	if mainTray == nil {
		mainTray = systray.New("OpenLinkHub system tray", "OpenLinkHub devices", icon.Data)
		go mainTray.Run()
	}
	for i := range *trays {
		entry := &(*trays)[i]
		if entry.serial == serial {
			entry.batteryItem.SetTitle(battString)
			entry.updated = true

			return nil
		}
	}

	batteryItem := mainTray.AddMenuItem(battString, "", nil)
	batteryItem.SetEnabled(false)

	*trays = append(*trays, traysEntry{
		serial:      serial,
		batteryItem: *batteryItem,
		updated:     true,
	})

	return nil
}

func cleanTrays(trays *[]traysEntry) {
	for i := 0; i < len(*trays); {
		if !(*trays)[i].updated {
			(*trays)[i].batteryItem.Quit()
			*trays = append((*trays)[:i], (*trays)[i+1:]...)
		} else {
			i++
		}
	}
}

func updateBattray(trays *[]traysEntry) {
	for i := range *trays {
		(*trays)[i].updated = false
	}

	icons, err := GetIcons()
	if err != nil {
		fmt.Printf("Error getting icons: %v", err)
	}

	stats, err := fetchDevices()
	if err != nil {
		fmt.Printf("Error getting devices: %v\n", err)
		return
	}

	for serial, device := range stats.Data {
		_ = handleDevice(serial, device, trays)
	}

	cleanTrays(trays)

	var deviceCount int = 0
	var totalLevel int = 0

	for _, device := range stats.Data {
		deviceCount++
		totalLevel += device.Level
	}

	totalLevel = totalLevel / deviceCount

	if totalLevel >= 75 {
		mainTray.SetIcon(icons.High)
	} else if totalLevel >= 25 {
		mainTray.SetIcon(icons.Normal)
	} else {
		mainTray.SetIcon(icons.Low)
	}
}

func StartBatteryModule() {
	trays := []traysEntry{}
	var traysMu sync.Mutex

	traysMu.Lock()
	updateBattray(&trays)
	traysMu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			traysMu.Lock()
			updateBattray(&trays)
			traysMu.Unlock()
		}
	}()
}
