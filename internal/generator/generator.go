package generator

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"one_c_swagger/internal/models"
	"one_c_swagger/internal/reader"
	"strings"
)

// updateRefsInContext recursively traverses an interface and updates all schema references.
func updateRefsInContext(data interface{}, serviceName string, localSchemaNames map[string]bool) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if key == "$ref" {
				if refStr, ok := val.(string); ok && strings.HasPrefix(refStr, "#/components/schemas/") {
					originalName := strings.TrimPrefix(refStr, "#/components/schemas/")
					if localSchemaNames[originalName] {
						v["$ref"] = fmt.Sprintf("#/components/schemas/%s_%s", serviceName, originalName)
					}
				}
			} else {
				updateRefsInContext(val, serviceName, localSchemaNames)
			}
		}
	case []interface{}:
		for _, item := range v {
			updateRefsInContext(item, serviceName, localSchemaNames)
		}
	}
}

func GenerateOpenAPI(services []reader.HTTPService, configs map[string]*reader.SwaggerConfig, allServicesConfig *reader.AllServicesConfig, log *slog.Logger) (*models.OpenAPI, error) {
	openapi := &models.OpenAPI{
		OpenAPI: "3.0.0",
		Info:    models.Info{Title: "1C HTTP Services", Version: "1.0.0"},
		Paths:   make(map[string]models.PathItem),
		Components: models.Components{
			Schemas:         make(map[string]interface{}),
			SecuritySchemes: make(map[string]models.SecurityScheme),
			Parameters:      make(map[string]models.Parameter),
			Headers:         make(map[string]models.Header),
			Responses:       make(map[string]interface{}),
		},
	}

	// --- PASS 0: Process Global Config ---
	globalSchemaNames := make(map[string]bool)
	if allServicesConfig != nil {
		openapi.Servers = allServicesConfig.Servers
		for k, v := range allServicesConfig.Components.SecuritySchemes {
			openapi.Components.SecuritySchemes[k] = v
		}
		for k, v := range allServicesConfig.Components.Parameters {
			openapi.Components.Parameters[k] = v
		}
		for k, v := range allServicesConfig.Components.Headers {
			openapi.Components.Headers[k] = v
		}
		for k, v := range allServicesConfig.Components.Responses {
			openapi.Components.Responses[k] = v
		}
		for k, v := range allServicesConfig.Components.Schemas {
			openapi.Components.Schemas[k] = v
			globalSchemaNames[k] = true
		}
	}

	// --- PASS 1: Collect and Rename Service-Specific Schemas ---
	for serviceName, config := range configs {
		if config.Components.Schemas != nil {
			for schemaName, schemaData := range config.Components.Schemas {
				newSchemaName := fmt.Sprintf("%s_%s", serviceName, schemaName)
				openapi.Components.Schemas[newSchemaName] = schemaData
			}
		}
	}

	// --- PASS 2: Update all $refs with context ---
	for serviceName, config := range configs {
		localSchemaNames := make(map[string]bool)
		if config.Components.Schemas != nil {
			for schemaName := range config.Components.Schemas {
				localSchemaNames[schemaName] = true
			}
		}
		updateRefsInContext(config.Paths, serviceName, localSchemaNames)
		updateRefsInContext(config.Components, serviceName, localSchemaNames)
	}

	// --- PASS 3: Process services and build paths ---
	for _, service := range services {
		openapi.Tags = append(openapi.Tags, models.Tag{Name: service.Properties.Name})
		swaggerConfig, hasSwaggerConfig := configs[service.Properties.Name]

		if hasSwaggerConfig {
			for name, scheme := range swaggerConfig.Components.SecuritySchemes {
				openapi.Components.SecuritySchemes[name] = scheme
			}
		}

		for _, urlTemplate := range service.URLTemplates {
			path := fmt.Sprintf("/%s/%s", strings.Trim(service.Properties.RootURL, "/"), strings.Trim(urlTemplate.Properties.Template, "/"))
			pathItem, ok := openapi.Paths[path]
			if !ok {
				pathItem = models.PathItem{}
			}

			for _, method := range urlTemplate.Methods {
				if strings.ToUpper(method.Properties.HTTPMethod) == "ANY" {
					pathItem.XAnyMethod = true
					continue
				}

				// 1. Get the overlay operation from the supplement file
				var overlayOp *models.Operation
				if hasSwaggerConfig && swaggerConfig.Paths != nil {
					if pathConfig, ok := swaggerConfig.Paths[urlTemplate.Properties.Template]; ok {
						pathConfigBytes, _ := json.Marshal(pathConfig)
						var pathItemConfig models.PathItem
						json.Unmarshal(pathConfigBytes, &pathItemConfig)
						switch strings.ToUpper(method.Properties.HTTPMethod) {
						case "GET":
							overlayOp = pathItemConfig.Get
						case "POST":
							overlayOp = pathItemConfig.Post
						}
					}
				}
				if overlayOp == nil {
					overlayOp = &models.Operation{}
				}

				// 2. Create the final operation and fill from 1C data
				finalOp := overlayOp
				if finalOp.Summary == "" {
					finalOp.Summary = method.Properties.Name
				}
				if finalOp.OperationID == "" {
					finalOp.OperationID = fmt.Sprintf("%s%s", service.Properties.Name, method.Properties.Name)
				}
				finalOp.Tags = []string{service.Properties.Name}

				// 3. Merge Responses
				if finalOp.Responses == nil {
					finalOp.Responses = make(models.Responses)
				}
				for code := range openapi.Components.Responses {
					if _, ok := finalOp.Responses[code]; !ok {
						finalOp.Responses[code] = map[string]string{"$ref": fmt.Sprintf("#/components/responses/%s", code)}
					}
				}

				// 4. Polish all non-ref responses with global headers
				for code, respIntf := range finalOp.Responses {
					respBytes, _ := json.Marshal(respIntf)
					var resp models.Response
					if json.Unmarshal(respBytes, &resp) == nil && !strings.Contains(string(respBytes), "$ref") {
						if resp.Headers == nil {
							resp.Headers = make(map[string]interface{})
						}
						for name := range openapi.Components.Headers {
							if _, ok := resp.Headers[name]; !ok {
								ref := fmt.Sprintf("#/components/headers/%s", name)
								resp.Headers[name] = map[string]string{"$ref": ref}
							}
						}
						finalOp.Responses[code] = resp
					}
				}

				// 5. Apply security
				if hasSwaggerConfig && len(swaggerConfig.Security) > 0 {
					finalOp.Security = swaggerConfig.Security
				} else if finalOp.Security == nil {
					finalOp.Security = []models.SecurityRequirement{{"basicAuth": {}}}
				}

				switch strings.ToUpper(method.Properties.HTTPMethod) {
				case "GET":
					pathItem.Get = finalOp
				case "POST":
					pathItem.Post = finalOp
				}
			}
			openapi.Paths[path] = pathItem
		}
	}

	return openapi, nil
}

func ToJSON(openapi *models.OpenAPI) (string, error) {
	jsonBytes, err := json.MarshalIndent(openapi, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}