package model

import (
	"strings"
)

type (
	Path struct {
		Tags       []string            `json:"tags,omitempty"`
		Summary    string              `json:"summary,omitempty"`
		Consumes   []string            `json:"consumes,omitempty"`
		Produces   []string            `json:"produces,omitempty"`
		Parameters []Parameter         `json:"parameters,omitempty"`
		Responses  map[string]Response `json:"responses,omitempty"`
	}

	Parameter struct {
		In          string            `json:"in,omitempty"`
		Name        string            `json:"name,omitempty"`
		Description string            `json:"description,omitempty"`
		Required    bool              `json:"required,omitempty"`
		Default     interface{}       `json:"default,omitempty"`
		Schema      map[string]string `json:"schema,omitempty"`
		Type        string            `json:"type,omitempty"`
		Format      string            `json:"format,omitempty"`
		Example     interface{}       `json:"x-example,omitempty"`
	}

	Response struct {
		Description string `json:"description,omitempty"`
		Schema      Schema `json:"schema,omitempty"`
	}

	Schema struct {
		Type   string            `json:"type,omitempty"`
		RefRaw string            `json:"$ref,omitempty"`
		Items  map[string]string `json:"items,omitempty"`
	}
)

func (p *Parameter) IsBody() bool {
	return strings.EqualFold(p.In, "body")
}

func (p *Parameter) TypeName() string {
	if p.IsBody() {
		return formatModelName(p.Schema["$ref"])
	}
	return p.Type
}

func formatModelName(s string) string {
	return strings.ReplaceAll(s, "#/definitions/", "")
}

func (schema Schema) IsVoid() bool {
	return schema.RefRaw != "" && schema.Items == nil
}

func (schema Schema) IsList() bool {
	return strings.EqualFold(schema.Type, "array")
}

func (schema Schema) Ref() string {
	if schema.Items == nil {
		return "ManagerResponse"
	}
	return formatModelName(schema.Items["$ref"])
}
