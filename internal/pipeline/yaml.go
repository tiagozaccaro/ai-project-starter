package pipeline

import "gopkg.in/yaml.v3"

// yamlUnmarshal parses YAML into a value.
func yamlUnmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// yamlMarshal serializes a value to YAML.
func yamlMarshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// isTerminal reports whether the given file is an interactive terminal.
func isTerminal(f *osFile) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & osModeCharDevice) != 0
}
