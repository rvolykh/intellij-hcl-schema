package intellij

type ProviderSchema struct {
	Name        string                `json:"name"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	Provider    SchemaInfo            `json:"provider"`
	Resources   map[string]SchemaInfo `json:"resources"`
	DataSources map[string]SchemaInfo `json:"data-sources"`
}

type SchemaInfo map[string]SchemaDefinition

type SchemaDefinition struct {
	Type        string `json:",omitempty"`
	Description string `json:",omitempty"`
	Deprecated  string `json:",omitempty"`
	Optional    bool   `json:",omitempty"`
	Required    bool   `json:",omitempty"`
	Computed    bool   `json:",omitempty"`
	Sensitive   bool   `json:",omitempty"`
}

// TODO: required support complex types, not yet implemented
type SchemaElement struct {
	// One of ValueType or "SchemaElements" or "SchemaInfo"
	Type string `json:",omitempty"`
	// Set for simple types (from ValueType)
	Value string `json:",omitempty"`
	// Set if Type == "SchemaElements"
	ElementsType string `json:",omitempty"`
	// Set if Type == "SchemaInfo"
	Info SchemaInfo `json:",omitempty"`
}
