package model

import (
	"github.com/virtuald/go-ordered-json"
	"strings"
)

type Definition struct {
	Title         string             `json:"title,omitempty"`
	Type          string             `json:"type,omitempty"`
	Required      []string           `json:"required,omitempty"`
	PropertiesRaw json.OrderedObject `json:"properties,omitempty"`
}

func (d *Definition) Properties() []Property {
	properties := make([]Property, 0, len(d.PropertiesRaw))
	for _, v := range d.PropertiesRaw {
		bytes, _ := json.Marshal(v.Value)
		var p Property
		err := json.Unmarshal(bytes, &p)
		if err != nil {
			panic(err)
		}
		p.Name = v.Key
		properties = append(properties, p)
	}

	return properties
}

type Property struct {
	Name        string
	Type        string      `json:"type,omitempty"`
	Format      string      `json:"format,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Description string      `json:"description,omitempty"`
	ReadOnly    bool        `json:"readOnly,omitempty"`
	RefRaw      string      `json:"$ref,omitempty"`
}

func (p *Property) Ref() string {
	return strings.ReplaceAll(p.RefRaw, "#/definitions/", "")
}
