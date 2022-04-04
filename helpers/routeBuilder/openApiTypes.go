package routeBuilder

// OpenAPI implements the root OpenAPI Response Object
// Ref: https://swagger.io/specification/#oas-object
type OpenAPI struct {
	OpenAPI    string                     `json:"openapi"` // Version
	Info       *OpenAPIInfo               `json:"info,omitempty"`
	Servers    []OpenAPIServer            `json:"servers,omitempty"`
	Paths      map[string]OpenAPIPathItem `json:"paths,omitempty"`
	Components *OpenAPIComponents         `json:"components,omitempty"`
	// Security   []OpenAPISecurityRequirement `json:"security,omitempty"` // TODO
	Tags []Tag `json:"tags,omitempty"`
	// ExternalDocs *OpenAPIExternalDocumentation `json:"externalDocs,omitempty"` // TODO
}

// OpenAPIComponents implements the OpenAPI Components
// Ref: https://swagger.io/specification/#components-object
type OpenAPIComponents struct {
	Schemas map[string]any `json:"schemas,omitempty"`
	// Responses map[string] `json:"responses,omitempty"`
	// Parameters map[string] `json:"parameters,omitempty"`
	// Examples map[string] `json:"examples,omitempty"`
	// RequestBodies map[string] `json:"requestBodies,omitempty"`
	// Headers map[string] `json:"headers,omitempty"`
	// SecuritySchemes map[string] `json:"securitySchemes,omitempty"`
	// Links map[string] `json:"links,omitempty"`
	// Callbacks map[string] `json:"callbacks,omitempty"`
}

// OpenAPIInfo contains information about the API.
// Ref: https://swagger.io/specification/#info-object
type OpenAPIInfo struct {
	Title          string         `json:"title,omitempty"`
	Description    string         `json:"description,omitempty"`
	TermsOfService string         `json:"termsOfService,omitempty"`
	Contact        OpenAPIContact `json:"contact,omitempty"`
	License        OpenAPILicense `json:"license,omitempty"`
	Version        string         `json:"version,omitempty"`
}

// OpenAPIContact implements the OpenAPI Contact Object
// Ref: https://swagger.io/specification/#contact-object
type OpenAPIContact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// OpenAPILicense implements the OpenAPI License Object
// Ref: https://swagger.io/specification/#license-object
type OpenAPILicense struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// OpenAPIPathItem implements the OpenAPI Path Item Object
// Ref: https://swagger.io/specification/#path-item-object
type OpenAPIPathItem struct {
	Ref         string             `json:"$ref,omitempty"`
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Get         *OpenAPIOperation  `json:"get,omitempty"`
	Put         *OpenAPIOperation  `json:"put,omitempty"`
	Post        *OpenAPIOperation  `json:"post,omitempty"`
	Delete      *OpenAPIOperation  `json:"delete,omitempty"`
	Options     *OpenAPIOperation  `json:"options,omitempty"`
	Head        *OpenAPIOperation  `json:"head,omitempty"`
	Patch       *OpenAPIOperation  `json:"patch,omitempty"`
	Trace       *OpenAPIOperation  `json:"trace,omitempty"`
	Servers     []OpenAPIServer    `json:"servers,omitempty"`
	Parameters  []OpenAPIParameter `json:"parameters,omitempty"`
}

// OpenAPIOperation implements the OpenAPI Operation Object
// Ref: https://swagger.io/specification/#operation-object
type OpenAPIOperation struct {
	Tags        []string `json:"tags,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	// ExternalDocs any `json:"externalDocs,omitempty"` // TODO
	OperationID string                     `json:"operationId,omitempty"`
	Parameters  []OpenAPIParameter         `json:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse `json:"responses,omitempty"`
	// Callbacks any `json:"callbacks,omitempty"` // TODO
	Deprecated bool `json:"deprecated,omitempty"`
	// Security   []OpenAPISecurityRequirement `json:"security,omitempty"` // TODO
	Servers []OpenAPIServer `json:"servers,omitempty"`
}

// OpenAPIRequestBody implements the OpenAPI Request Body Object
// Ref: https://swagger.io/specification/#request-body-object
type OpenAPIRequestBody struct {
	Description string                      `json:"description,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty"`
	Required    bool                        `json:"required,omitempty"`
}

// OpenAPIResponse implements the OpenAPI Response Object
// Ref: https://swagger.io/specification/#response-object
type OpenAPIResponse struct {
	Description string `json:"description,omitempty"`
	// Headers map[string]OpenAPIHeader // TODO
	Content map[string]OpenAPIMediaType `json:"content,omitempty"`
	// Links map[string]OpenAPILink // TODO
}

// OpenAPIServer implements the OpenAPI Server Object
// Ref: https://swagger.io/specification/#server-object
type OpenAPIServer struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	// Variables   map[string]OpenAPIServerVariable `json:"variables,omitempty"` // TODO
}

// OpenAPIParameter implements the OpenAPI Parameter Object
// Ref: https://swagger.io/specification/#parameter-object
type OpenAPIParameter struct {
	Name            string `json:"name,omitempty"`
	In              string `json:"in,omitempty"`
	Description     string `json:"description,omitempty"`
	Required        bool   `json:"required,omitempty"`
	Deprecated      bool   `json:"deprecated,omitempty"`
	AllowEmptyValue bool   `json:"allowEmptyValue,omitempty"`
	//
	Style         string `json:"style,omitempty"`
	Explode       bool   `json:"explode,omitempty"`
	AllowReserved bool   `json:"allowReserved,omitempty"`
	Schema        any    `json:"schema,omitempty"`
	Example       any    `json:"example,omitempty"`
	Examples      any    `json:"examples,omitempty"`
}

// OpenAPIMediaType implements the me OpenAPI Media Type Object
// Ref: https://swagger.io/specification/#media-type-object
type OpenAPIMediaType struct {
	Schema   any            `json:"schema,omitempty"`
	Example  any            `json:"example,omitempty"`
	Examples map[string]any `json:"examples,omitempty"`
	// Encoding map[string]OpenAPIEncoding `json:"encoding,omitempty"` // TODO
}
