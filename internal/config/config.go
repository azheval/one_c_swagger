package config

import (
	"encoding/json"
	"os"
)

// Config структура для хранения настроек приложения
type Config struct {
	Log     Log     `json:"log"`
	Project Project `json:"project"`
}

// Log структура для хранения настроек логирования
type Log struct {
	LogPath  string `json:"log_path"`
	LogLevel string `json:"log_level"`
}

// Project структура для хранения настроек проекта
type Project struct {
	ConfigurationPath         string `json:"configuration_path"`
	ExtensionsPath           string `json:"extensions_path"`
	OutPath                   string `json:"out_path"`
	SwaggerConfigPath         string `json:"swagger_config_path"`
	AllServicesConfigFileName string `json:"all_services_config_filename"`
}

// LoadConfig читает и разбирает файл конфигурации
func LoadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
