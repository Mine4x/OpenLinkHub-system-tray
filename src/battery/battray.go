package battery

import (
	"fmt"
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
	println("Handling device")

	battString := fmt.Sprintf("🔋 %s: %d%%", device.Device, device.Level)

	icons, err := GetIcons()
	if err != nil {
		return fmt.Errorf("Couldn't get icons: %v", err)
	}

	for _, entry := range *trays {
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

	newTray := systray.New(fmt.Sprintf("🔋 %s: %d%%", device.Device, device.Level), "OpenLinkHub device", icons.Normal)
	batteryItem := newTray.AddMenuItem(battString, "", nil)

	batteryItem.SetEnabled(false)

	newEntry := traysEntry{
		serial:      serial,
		tray:        *newTray,
		batteryItem: *batteryItem,
		updated:     true,
	}

	*trays = append(*trays, newEntry)

	fmt.Println("Running tray")
	newTray.Run() // FIXME: Running tray on same Dbus as existing tray

	return nil
}

func cleanTrays(trays *[]traysEntry) {
	for i, entry := range *trays {
		if entry.updated != true {
			entry.tray.Quit()
			*trays = append((*trays)[:i], (*trays)[i+1:]...)
			cleanTrays(trays)
			return
		}
	}
}

func updateBattray(trays *[]traysEntry) {
	for _, entry := range *trays {
		entry.updated = false
	}

	stats, err := fetchDevices()
	if err != nil {
		fmt.Printf("Error getting devices: %v", err)
		return
	}

	for serial, device := range stats.Data {
		handleDevice(serial, device, trays)
	}

	cleanTrays(trays)
}

func StartBatteryModule() {
	trays := []traysEntry{}

	updateBattray(&trays)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateBattray(&trays)
		}
	}()
}
