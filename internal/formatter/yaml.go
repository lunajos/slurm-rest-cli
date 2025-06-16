package formatter

import (
	"gopkg.in/yaml.v3"
)

// YAMLFormatter formats data as YAML.
type YAMLFormatter struct{}

// Format formats the data as YAML.
func (f *YAMLFormatter) Format(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
