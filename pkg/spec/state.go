package spec

// State handler spec.
type State struct {
	Set    map[string]string `yaml:"set"`
	Delete []string          `yaml:"delete"`
}
