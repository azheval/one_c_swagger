package models

type OpenAPI struct {
	OpenAPI    string               `json:"openapi"`
	Info       Info                 `json:"info"`
	Servers    []Server             `json:"servers,omitempty"`
	Tags       []Tag                `json:"tags,omitempty"`
	Paths      map[string]PathItem  `json:"paths"`
	Components Components           `json:"components"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type PathItem struct {
	Get       *Operation `json:"get,omitempty"`
	Post      *Operation `json:"post,omitempty"`
	Put       *Operation `json:"put,omitempty"`
	Delete    *Operation `json:"delete,omitempty"`
	Head      *Operation `json:"head,omitempty"`
	Patch     *Operation `json:"patch,omitempty"`
	Merge     *Operation `json:"merge,omitempty"`
	Options   *Operation `json:"options,omitempty"`
	Trace     *Operation `json:"trace,omitempty"`
	Connect   *Operation `json:"connect,omitempty"`
	Propfind  *Operation `json:"propfind,omitempty"`
	Proppatch *Operation `json:"proppatch,omitempty"`
	Move      *Operation `json:"move,omitempty"`
	Copy      *Operation `json:"copy,omitempty"`
	Lock      *Operation `json:"lock,omitempty"`
	Unlock    *Operation `json:"unlock,omitempty"`
	Mkcol     *Operation `json:"mkcol,omitempty"`
	XAnyMethod bool      `json:"x-any-method,omitempty"`
}

type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary"`
	OperationID string                `json:"operationId"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody interface{}           `json:"requestBody,omitempty"`
	Responses   Responses             `json:"responses"`
	Security    []SecurityRequirement `json:"security,omitempty"`
	Servers     []Server              `json:"servers,omitempty"`
}

type Responses map[string]interface{}

type Response struct {
	Description string                 `json:"description"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
	Content     map[string]MediaType   `json:"content,omitempty"`
}

type Header struct {
	Description string    `json:"description,omitempty"`
	Schema      SchemaRef `json:"schema"`
}

type MediaType struct {
	Schema  interface{} `json:"schema,omitempty"`
	Example interface{} `json:"example,omitempty"`
}

type Components struct {
	Schemas         map[string]interface{}    `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty"`
	Responses       map[string]interface{}    `json:"responses,omitempty"`
}

type SecurityScheme struct {
	Type         string `json:"type"`
	Description  string `json:"description,omitempty"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

type SecurityRequirement map[string][]string

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      *SchemaRef  `json:"schema,omitempty"`
}

type SchemaRef struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}
