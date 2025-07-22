package helpers

// Config represents the structure of the YAML config file
type navConfig struct {
	Folders  []string `yaml:"Folders"`
	MaxDepth int      `yaml:"maxDepth"`
	Comments string   `yaml:"comments"`
}
