package reader

import (
	"encoding/xml"
	"one_c_swagger/internal/models"
)

type MetaDataObject struct {
	XMLName     xml.Name    `xml:"MetaDataObject"`
	HTTPService HTTPService `xml:"HTTPService"`
}

type HTTPService struct {
	XMLName      xml.Name              `xml:"HTTPService"`
	UUID         string                `xml:"uuid,attr"`
	Properties   HTTPServiceProperties `xml:"Properties"`
	URLTemplates []URLTemplate         `xml:"ChildObjects>URLTemplate"`
}

type HTTPServiceProperties struct {
	Name    string `xml:"Name"`
	Synonym struct {
		Item struct {
			Lang    string `xml:"lang"`
			Content string `xml:"content"`
		} `xml:"item"`
	} `xml:"Synonym"`
	RootURL string `xml:"RootURL"`
}

type URLTemplate struct {
	XMLName    xml.Name              `xml:"URLTemplate"`
	UUID       string                `xml:"uuid,attr"`
	Properties URLTemplateProperties `xml:"Properties"`
	Methods    []Method              `xml:"ChildObjects>Method"`
}

type URLTemplateProperties struct {
	Name    string `xml:"Name"`
	Synonym struct {
		Item struct {
			Lang    string `xml:"lang"`
			Content string `xml:"content"`
		} `xml:"item"`
	} `xml:"Synonym"`
	Template string `xml:"Template"`
}

type Method struct {
	XMLName    xml.Name         `xml:"Method"`
	UUID       string           `xml:"uuid,attr"`
	Properties MethodProperties `xml:"Properties"`
}

type MethodProperties struct {
	Name       string `xml:"Name"`
	Synonym    string `xml:"Synonym"`
	HTTPMethod string `xml:"HTTPMethod"`
	Handler    string `xml:"Handler"`
}

type SwaggerConfig struct {
	Servers    []models.Server              `json:"servers,omitempty"`
	Components models.Components            `json:"components,omitempty"`
	Security   []models.SecurityRequirement `json:"security,omitempty"`
	Paths      map[string]interface{}       `json:"paths"`
}

type AllServicesConfig struct {
	Servers    []models.Server   `json:"servers,omitempty"`
	Components models.Components `json:"components,omitempty"`
}
