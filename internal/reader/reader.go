package reader

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	utf8BOM = []byte{0xEF, 0xBB, 0xBF}
)

func ReadHTTPServices(path string, log *slog.Logger) ([]HTTPService, error) {
	log.Info("Reading http services", "path", path)
	var services []HTTPService
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error accessing path", "path", path, "error", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xml") {
			log.Info("Found http service file", "path", path)
			content, err := os.ReadFile(path)
			if err != nil {
				log.Error("Error reading file", "path", path, "error", err)
				return err
			}

			// Check for UTF-8 BOM
			if bytes.HasPrefix(content, utf8BOM) {
				content = bytes.TrimPrefix(content, utf8BOM)
			}

			var data MetaDataObject
			if err := xml.Unmarshal(content, &data); err != nil {
				log.Error("Error decoding xml", "path", path, "error", err)
				return err
			}

			log.Info("Successfully parsed http service", "path", path, "service", data.HTTPService.Properties.Name)
			log.Debug("Parsed data", "data", data)
			services = append(services, data.HTTPService)

		}
		return nil
	})
	if err != nil {
		log.Error("Error walking through http services directory", "path", path, "error", err)
		return nil, err
	}
	return services, nil
}

func ReadSwaggerConfigFile(swaggerConfigPath, serviceName string, log *slog.Logger) (*SwaggerConfig, error) {
	configPath := filepath.Join(swaggerConfigPath, serviceName+".json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	log.Info("Found swagger-config file", "path", configPath)
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config SwaggerConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ReadAllServicesConfigFile(swaggerConfigPath, fileName string, log *slog.Logger) (*AllServicesConfig, error) {
	configPath := filepath.Join(swaggerConfigPath, fileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	log.Info("Found all services config file", "path", configPath)
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AllServicesConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
