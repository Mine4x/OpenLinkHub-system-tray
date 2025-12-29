package battery

import (
	"fmt"
	"sync"
	"time"

	"github.com/Mine4x/OpenLinkHub-system-tray/src/systray"
)

type traysEntry struct {
	serial      string
	tray        systray.Tray
	batteryItem systray.MenuItem
	updated     bool
}

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

	icons, err := GetIcons()
	if err != nil {
		return err
	}

	for i := range *trays {
		entry := &(*trays)[i]
		if entry.serial == serial {
			if device.Level >= 80 {
				entry.tray.SetIcon(icons.High)
			} else if device.Level >= 35 {
				entry.tray.SetIcon(icons.Normal)
			} else {
				entry.tray.SetIcon(icons.Low)
			}
			entry.tray.SetTitle(battString)
			entry.batteryItem.SetTitle(battString)
			entry.updated = true
			return nil
		}
	}

	newTray := systray.New(battString, "OpenLinkHub device", icons.Normal)
	batteryItem := newTray.AddMenuItem(battString, "", nil)
	batteryItem.SetEnabled(false)

	*trays = append(*trays, traysEntry{
		serial:      serial,
		tray:        *newTray,
		batteryItem: *batteryItem,
		updated:     true,
	})

	go newTray.Run()

	return nil
}

func cleanTrays(trays *[]traysEntry) {
	for i := 0; i < len(*trays); {
		if !(*trays)[i].updated {
			(*trays)[i].tray.Quit()
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

	stats, err := fetchDevices()
	if err != nil {
		fmt.Printf("Error getting devices: %v\n", err)
		return
	}

	for serial, device := range stats.Data {
		_ = handleDevice(serial, device, trays)
	}

	cleanTrays(trays)
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
