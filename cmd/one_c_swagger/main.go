package main

import (
	"flag"
	"fmt"
	"log"
	"one_c_swagger/internal/config"
	"one_c_swagger/internal/generator"
	"one_c_swagger/internal/logger"
	"one_c_swagger/internal/merger"
	"one_c_swagger/internal/reader"
	"os"
	"path/filepath"
)

var (
	version = "dev"
	build   = "local"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and build information")
	configPath := flag.String("config", "configs/config.json", "Path to the configuration file")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\nBuild: %s\n", version, build)
		return
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	slog := logger.New(cfg.Log.LogLevel, cfg.Log.LogPath)

	slog.Info("Starting one_c_swagger", "version", version, "build", build)
	slog.Info("Config loaded", "config", cfg)

	var baseServices []reader.HTTPService
	var extServices []reader.HTTPService

	// Read base configuration
	if cfg.Project.ConfigurationPath != "" {
		if _, err := os.Stat(cfg.Project.ConfigurationPath); !os.IsNotExist(err) {
			httpServicesPath := filepath.Join(cfg.Project.ConfigurationPath, "HTTPServices")
			if _, err := os.Stat(httpServicesPath); !os.IsNotExist(err) {
				services, err := reader.ReadHTTPServices(httpServicesPath, slog)
				if err != nil {
					slog.Error("Error reading http services from configuration", "path", httpServicesPath, "error", err)
				} else {
					slog.Info("Found services in configuration", "count", len(services))
					baseServices = append(baseServices, services...)
				}
			}
		} else {
			slog.Warn("Configuration path does not exist", "path", cfg.Project.ConfigurationPath)
		}
	}

	// Read extensions
	if cfg.Project.ExtensionsPath != "" {
		if _, err := os.Stat(cfg.Project.ExtensionsPath); !os.IsNotExist(err) {
			extensions, err := os.ReadDir(cfg.Project.ExtensionsPath)
			if err != nil {
				slog.Error("Error reading extensions directory", "path", cfg.Project.ExtensionsPath, "error", err)
			} else {
				for _, ext := range extensions {
					if ext.IsDir() {
						extHttpServicesPath := filepath.Join(cfg.Project.ExtensionsPath, ext.Name(), "HTTPServices")
						if _, err := os.Stat(extHttpServicesPath); !os.IsNotExist(err) {
							services, err := reader.ReadHTTPServices(extHttpServicesPath, slog)
							if err != nil {
								slog.Error("Error reading http services from extension", "path", extHttpServicesPath, "error", err)
							} else {
								slog.Info("Found services in extension", "extension", ext.Name(), "count", len(services))
								extServices = append(extServices, services...)
							}
						}
					}
				}
			}
		} else {
			slog.Warn("Extensions path does not exist", "path", cfg.Project.ExtensionsPath)
		}
	}

	mergedServices := merger.MergeServices(baseServices, extServices)
	slog.Info("Total http services after merge", "count", len(mergedServices))

	// Read all services config
	var allServicesConfig *reader.AllServicesConfig
	if cfg.Project.SwaggerConfigPath != "" && cfg.Project.AllServicesConfigFileName != "" {
		allServicesConfig, err = reader.ReadAllServicesConfigFile(cfg.Project.SwaggerConfigPath, cfg.Project.AllServicesConfigFileName, slog)
		if err != nil {
			slog.Error("Error reading all services config file", "error", err)
		}
	}

	// Read swagger configs
	swaggerConfigs := make(map[string]*reader.SwaggerConfig)
	if cfg.Project.SwaggerConfigPath != "" {
		for _, service := range mergedServices {
			swaggerConfig, err := reader.ReadSwaggerConfigFile(cfg.Project.SwaggerConfigPath, service.Properties.Name, slog)
			if err != nil {
				slog.Error("Error reading swagger-config.json", "service", service.Properties.Name, "error", err)
			}
			if swaggerConfig != nil {
				swaggerConfigs[service.Properties.Name] = swaggerConfig
			}
		}
	}

	// Generate OpenAPI spec
	openapi, err := generator.GenerateOpenAPI(mergedServices, swaggerConfigs, allServicesConfig, slog)
	if err != nil {
		slog.Error("Error generating OpenAPI object", "error", err)
		return
	}

	// Create out directory if it doesn't exist
	if _, err := os.Stat(cfg.Project.OutPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.Project.OutPath, 0755); err != nil {
			slog.Error("Error creating output directory", "path", cfg.Project.OutPath, "error", err)
			return
		}
	}

	// Generate and save JSON spec
	jsonSpec, err := generator.ToJSON(openapi)
	if err != nil {
		slog.Error("Error generating JSON spec", "error", err)
	} else {
		jsonFile := filepath.Join(cfg.Project.OutPath, "openapi.json")
		if err := os.WriteFile(jsonFile, []byte(jsonSpec), 0644); err != nil {
			slog.Error("Error writing json spec file", "path", jsonFile, "error", err)
		} else {
			slog.Info("Successfully generated openapi.json", "path", jsonFile)
		}
	}

}