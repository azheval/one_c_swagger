package merger

import "one_c_swagger/internal/reader"

func MergeServices(baseServices, extServices []reader.HTTPService) []reader.HTTPService {
	serviceMap := make(map[string]reader.HTTPService)

	// Add base services to the map
	for _, service := range baseServices {
		serviceMap[service.Properties.Name] = service
	}

	// Merge extension services
	for _, extService := range extServices {
		if existingService, ok := serviceMap[extService.Properties.Name]; ok {
			// Service exists, merge URL templates
			mergedURLTemplates := mergeURLTemplates(existingService.URLTemplates, extService.URLTemplates)
			existingService.URLTemplates = mergedURLTemplates
			serviceMap[extService.Properties.Name] = existingService
		} else {
			// Service does not exist, add it
			serviceMap[extService.Properties.Name] = extService
		}
	}

	// Convert map back to slice
	var mergedServices []reader.HTTPService
	for _, service := range serviceMap {
		mergedServices = append(mergedServices, service)
	}

	return mergedServices
}

func mergeURLTemplates(baseTemplates, extTemplates []reader.URLTemplate) []reader.URLTemplate {
	templateMap := make(map[string]reader.URLTemplate)

	// Add base templates to the map
	for _, template := range baseTemplates {
		templateMap[template.Properties.Template] = template
	}

	// Merge extension templates
	for _, extTemplate := range extTemplates {
		if existingTemplate, ok := templateMap[extTemplate.Properties.Template]; ok {
			// Template exists, merge methods
			mergedMethods := mergeMethods(existingTemplate.Methods, extTemplate.Methods)
			existingTemplate.Methods = mergedMethods
			templateMap[extTemplate.Properties.Template] = existingTemplate
		} else {
			// Template does not exist, add it
			templateMap[extTemplate.Properties.Template] = extTemplate
		}
	}

	// Convert map back to slice
	var mergedTemplates []reader.URLTemplate
	for _, template := range templateMap {
		mergedTemplates = append(mergedTemplates, template)
	}

	return mergedTemplates
}

func mergeMethods(baseMethods, extMethods []reader.Method) []reader.Method {
	methodMap := make(map[string]reader.Method)

	// Add base methods to the map
	for _, method := range baseMethods {
		methodMap[method.Properties.Name] = method
	}

	// Merge extension methods
	for _, extMethod := range extMethods {
		// Always overwrite, as per "дополнять или переопределять"
		methodMap[extMethod.Properties.Name] = extMethod
	}

	// Convert map back to slice
	var mergedMethods []reader.Method
	for _, method := range methodMap {
		mergedMethods = append(mergedMethods, method)
	}

	return mergedMethods
}
