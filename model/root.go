package model

type (
	Root struct {
		BasePath    string                     `json:"basePath,omitempty"`
		Tags        []Tag                      `json:"tags,omitempty"`
		Paths       map[string]map[string]Path `json:"paths,omitempty"`
		Definitions map[string]Definition      `json:"definitions,omitempty"`
	}

	Tag struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}
)
