// Package openapi2 parses and writes OpenAPI 2 specifications.
//
// Does not cover all elements of OpenAPI 2.
// When OpenAPI version 3 is backwards-compatible with version 2, version 3 elements have been used.
//
// The specification:
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md
package openapi2

import (
	"fmt"

	"github.com/mbilski/kin-openapi/jsoninfo"
	"github.com/mbilski/kin-openapi/openapi3"
)

type Swagger struct {
	Info                openapi3.Info                  `json:"info"`
	ExternalDocs        *openapi3.ExternalDocs         `json:"externalDocs,omitempty"`
	Schemes             []string                       `json:"schemes,omitempty"`
	Host                string                         `json:"host,omitempty"`
	BasePath            string                         `json:"basePath,omitempty"`
	Paths               map[string]*PathItem           `json:"paths,omitempty"`
	Definitions         map[string]*openapi3.SchemaRef `json:"definitions,omitempty,noref"`
	Parameters          map[string]*Parameter          `json:"parameters,omitempty,noref"`
	Responses           map[string]*Response           `json:"responses,omitempty,noref"`
	SecurityDefinitions map[string]*SecurityScheme     `json:"securityDefinitions,omitempty"`
	Security            SecurityRequirements           `json:"security,omitempty"`
	Tags                openapi3.Tags                  `json:"tags,omitempty"`
}

func (swagger *Swagger) AddOperation(path string, method string, operation *Operation) {
	paths := swagger.Paths
	if paths == nil {
		paths = make(map[string]*PathItem, 8)
		swagger.Paths = paths
	}
	pathItem := paths[path]
	if pathItem == nil {
		pathItem = &PathItem{}
		paths[path] = pathItem
	}
	pathItem.SetOperation(method, operation)
}

type PathItem struct {
	Ref        string     `json:"$ref,omitempty"`
	Delete     *Operation `json:"delete,omitempty"`
	Get        *Operation `json:"get,omitempty"`
	Head       *Operation `json:"head,omitempty"`
	Options    *Operation `json:"options,omitempty"`
	Patch      *Operation `json:"patch,omitempty"`
	Post       *Operation `json:"post,omitempty"`
	Put        *Operation `json:"put,omitempty"`
	Parameters Parameters `json:"parameters,omitempty"`
}

func (pathItem *PathItem) Operations() map[string]*Operation {
	operations := make(map[string]*Operation, 8)
	if v := pathItem.Delete; v != nil {
		operations["DELETE"] = v
	}
	if v := pathItem.Get; v != nil {
		operations["GET"] = v
	}
	if v := pathItem.Head; v != nil {
		operations["HEAD"] = v
	}
	if v := pathItem.Options; v != nil {
		operations["OPTIONS"] = v
	}
	if v := pathItem.Patch; v != nil {
		operations["PATCH"] = v
	}
	if v := pathItem.Post; v != nil {
		operations["POST"] = v
	}
	if v := pathItem.Put; v != nil {
		operations["PUT"] = v
	}
	return operations
}

func (pathItem *PathItem) GetOperation(method string) *Operation {
	switch method {
	case "DELETE":
		return pathItem.Delete
	case "GET":
		return pathItem.Get
	case "HEAD":
		return pathItem.Head
	case "OPTIONS":
		return pathItem.Options
	case "PATCH":
		return pathItem.Patch
	case "POST":
		return pathItem.Post
	case "PUT":
		return pathItem.Put
	default:
		panic(fmt.Errorf("Unsupported HTTP method '%s'", method))
	}
}

func (pathItem *PathItem) SetOperation(method string, operation *Operation) {
	switch method {
	case "DELETE":
		pathItem.Delete = operation
	case "GET":
		pathItem.Get = operation
	case "HEAD":
		pathItem.Head = operation
	case "OPTIONS":
		pathItem.Options = operation
	case "PATCH":
		pathItem.Patch = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	default:
		panic(fmt.Errorf("Unsupported HTTP method '%s'", method))
	}
}

type Operation struct {
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	ExternalDocs *openapi3.ExternalDocs `json:"externalDocs,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	OperationID  string                 `json:"operationId,omitempty"`
	Parameters   Parameters             `json:"parameters,omitempty"`
	Responses    map[string]*Response   `json:"responses"`
	Consumes     []string               `json:"consumes,omitempty"`
	Produces     []string               `json:"produces,omitempty"`
	Security     *SecurityRequirements  `json:"security,omitempty"`
}

type Parameters []*Parameter

type Parameter struct {
	Ref          string              `json:"$ref,omitempty"`
	In           string              `json:"in,omitempty"`
	Name         string              `json:"name,omitempty"`
	Description  string              `json:"description,omitempty"`
	Required     bool                `json:"required,omitempty"`
	UniqueItems  bool                `json:"uniqueItems,omitempty"`
	ExclusiveMin bool                `json:"exclusiveMinimum,omitempty"`
	ExclusiveMax bool                `json:"exclusiveMaximum,omitempty"`
	Schema       *openapi3.SchemaRef `json:"schema,omitempty"`
	Type         string              `json:"type,omitempty"`
	Format       string              `json:"format,omitempty"`
	Enum         []interface{}       `json:"enum,omitempty"`
	Minimum      *float64            `json:"minimum,omitempty"`
	Maximum      *float64            `json:"maximum,omitempty"`
	MinLength    uint64              `json:"minLength,omitempty"`
	MaxLength    *uint64             `json:"maxLength,omitempty"`
	Pattern      string              `json:"pattern,omitempty"`
	Items        *openapi3.SchemaRef `json:"items,omitempty"`
	MinItems     uint64              `json:"minItems,omitempty"`
	MaxItems     *uint64             `json:"maxItems,omitempty"`
	Default      interface{}         `json:"default,omitempty"`
}

type Response struct {
	Ref         string                 `json:"$ref,omitempty"`
	Description string                 `json:"description,omitempty"`
	Schema      *openapi3.SchemaRef    `json:"schema,omitempty"`
	Headers     map[string]*Header     `json:"headers,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
}

type Header struct {
	Ref         string `json:"$ref,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type SecurityRequirements []map[string][]string

type SecurityScheme struct {
	ExtensionProps

	Ref              string            `json:"$ref,omitempty"`
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type,omitempty"`
	In               string            `json:"in,omitempty"`
	Name             string            `json:"name,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
	Tags             openapi3.Tags     `json:"tags,omitempty"`
}

func (ss *SecurityScheme) MarshalJSON() ([]byte, error) {
	return jsoninfo.MarshalStrictStruct(ss)
}

func (ss *SecurityScheme) UnmarshalJSON(data []byte) error {
	return jsoninfo.UnmarshalStrictStruct(data, ss)
}

type ExtensionProps struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`
}

// Assert that the type implements the interface
var _ jsoninfo.StrictStruct = &ExtensionProps{}

// EncodeWith will be invoked by package "jsoninfo"
func (props *ExtensionProps) EncodeWith(encoder *jsoninfo.ObjectEncoder, value interface{}) error {
	for k, v := range props.Extensions {
		if err := encoder.EncodeExtension(k, v); err != nil {
			return err
		}
	}
	return encoder.EncodeStructFieldsAndExtensions(value)
}

// DecodeWith will be invoked by package "jsoninfo"
func (props *ExtensionProps) DecodeWith(decoder *jsoninfo.ObjectDecoder, value interface{}) error {
	if err := decoder.DecodeStructFieldsAndExtensions(value); err != nil {
		return err
	}
	source := decoder.DecodeExtensionMap()
	result := make(map[string]interface{}, len(source))
	for k, v := range source {
		result[k] = v
	}
	props.Extensions = result
	return nil
}
