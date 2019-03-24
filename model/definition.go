package model

type Definition struct {
	Title      string   `json:"title,omitempty"`
	Type       string   `json:"type,omitempty"`
	Required   []string `json:"required,omitempty"`
	Properties map[string]struct {
		Type        string      `json:"type,omitempty"`
		Format      string      `json:"format,omitempty"`
		Example     interface{} `json:"example,omitempty"`
		Description string      `json:"description,omitempty"`
		ReadOnly    bool        `json:"readOnly,omitempty"`
	} `json:"properties,omitempty"`
}
