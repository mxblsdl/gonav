package helpers

// Config represents the structure of the YAML config file
type navConfig struct {
	DefaultFolders []string `yaml:"defaultFolders"`
	MaxDepth       int                  `yaml:"maxDepth"`
}