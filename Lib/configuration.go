package lib

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Configuration holds server-side settings; never pass it to browser templates.
type Configuration struct {
	OBS struct {
		IP       string `json:"ip"`
		Port     int    `json:"port"`
		Password string `json:"password"`
	} `json:"obs"`
	Server struct {
		ListeningIP string `json:"listening_ip"`
		Port        int    `json:"port"`
	} `json:"server"`
}

// ReadConfiguration reads settings without exposing file contents in errors.
func ReadConfiguration(path string) (Configuration, error) {
	var config Configuration
	data, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("configuration.json okunamadı: %w", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return Configuration{}, fmt.Errorf("configuration.json geçerli yapılandırma JSON'u olmalıdır")
	}
	if net.ParseIP(config.OBS.IP) == nil || net.ParseIP(config.Server.ListeningIP) == nil {
		return Configuration{}, fmt.Errorf("configuration.json: obs.ip ve server.listening_ip geçerli IP adresleri olmalıdır")
	}
	if config.OBS.Port < 1 || config.OBS.Port > 65535 || config.Server.Port < 1 || config.Server.Port > 65535 {
		return Configuration{}, fmt.Errorf("configuration.json: portlar 1–65535 aralığında olmalıdır")
	}
	return config, nil
}
