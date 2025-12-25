package battery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Mine4x/OpenLinkHub-system-tray/src/config"
)

type BatteryResponse struct {
	Code   int                      `json:"code"`
	Status int                      `json:"status"`
	Data   map[string]BatteryDevice `json:"data"`
}

type BatteryDevice struct {
	Device     string `json:"Device"`
	Level      int    `json:"Level"`
	DeviceType int    `json:"DeviceType"`
}

type BatteryIcons struct {
	High   []byte `json:"High"`
	Normal []byte `json:"Normal"`
	Low    []byte `json:"Low"`
}

func getApiURL() (*string, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/batteryStats", conf.APIURL)

	return &apiURL, nil
}

func getIconPath() (*string, error) {
	conf, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &conf.IconsPath, nil
}

func GetBatteryStats() (*BatteryResponse, error) {
	apiURL, err := getApiURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get apiURL form config: %w", err)
	}

	resp, err := http.Get(*apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var batteryResp BatteryResponse
	if err := json.Unmarshal(body, &batteryResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &batteryResp, nil
}

func PrintBatteryStats() error {
	stats, err := GetBatteryStats()
	if err != nil {
		return err
	}

	fmt.Printf("API Response Code: %d\n", stats.Code)
	fmt.Printf("Status: %d\n\n", stats.Status)

	if len(stats.Data) == 0 {
		fmt.Println("No devices found")
		return nil
	}

	fmt.Println("Devices:")
	for serial, device := range stats.Data {
		fmt.Printf("  Serial: %s\n", serial)
		fmt.Printf("    Name: %s\n", device.Device)
		fmt.Printf("    Battery Level: %d%%\n", device.Level)
		fmt.Printf("    Device Type: %d\n", device.DeviceType)
		fmt.Println()
	}

	return nil
}

func GetDeviceBySerial(serial string) (*BatteryDevice, error) {
	stats, err := GetBatteryStats()
	if err != nil {
		return nil, err
	}

	device, exists := stats.Data[serial]
	if !exists {
		return nil, fmt.Errorf("device with serial %s not found", serial)
	}

	return &device, nil
}

func GetAllDevices() ([]BatteryDevice, error) {
	stats, err := GetBatteryStats()
	if err != nil {
		return nil, err
	}

	devices := make([]BatteryDevice, 0, len(stats.Data))
	for _, device := range stats.Data {
		devices = append(devices, device)
	}

	return devices, nil
}

func GetLowestBattery() (*BatteryDevice, string, error) {
	stats, err := GetBatteryStats()
	if err != nil {
		return nil, "", err
	}

	if len(stats.Data) == 0 {
		return nil, "", fmt.Errorf("no devices found")
	}

	var lowestDevice *BatteryDevice
	var lowestSerial string
	lowestLevel := 101 // Start above 100%

	for serial, device := range stats.Data {
		if device.Level < lowestLevel {
			lowestLevel = device.Level
			deviceCopy := device
			lowestDevice = &deviceCopy
			lowestSerial = serial
		}
	}

	return lowestDevice, lowestSerial, nil
}

func GetIcons() (*BatteryIcons, error) {
	Path, err := getIconPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get icon path: %w", err)
	}

	HighIcon, err := os.ReadFile(fmt.Sprintf("%s/battery_high.png", *Path))
	if err != nil {
		return nil, fmt.Errorf("failed to get icon path: %w", err)
	}
	NormalIcon, err := os.ReadFile(fmt.Sprintf("%s/battery_normal.png", *Path))
	if err != nil {
		return nil, fmt.Errorf("failed to get icon path: %w", err)
	}
	LowIcon, err := os.ReadFile(fmt.Sprintf("%s/battery_low.png", *Path))
	if err != nil {
		return nil, fmt.Errorf("failed to get icon path: %w", err)
	}

	var batteryIcons BatteryIcons
	batteryIcons.High = HighIcon
	batteryIcons.Normal = NormalIcon
	batteryIcons.Low = LowIcon

	return &batteryIcons, nil
}
