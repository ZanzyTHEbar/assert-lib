package assert

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Formatter interface for structured output
type Formatter interface {
	Format(assertData map[string]interface{}, stack string) string
}

// TextFormatter is the default plain text output format
type TextFormatter struct{}

func (f *TextFormatter) Format(assertData map[string]interface{}, stack string) string {
	output := "ASSERT\n"
	for key, value := range assertData {
		output += fmt.Sprintf("   %s=%v\n", key, value)
	}
	if stack != "" {
		output += fmt.Sprintf("%s\n", stack)
	}
	return output
}

// JSONFormatter for JSON output
type JSONFormatter struct{}

func (f *JSONFormatter) Format(assertData map[string]interface{}, stack string) string {
	data := map[string]interface{}{
		"assertData": assertData,
	}
	if stack != "" {
		data["stack"] = stack
	}
	out, _ := json.MarshalIndent(data, "", "  ")
	return string(out)
}

// YAMLFormatter for YAML output
type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(assertData map[string]interface{}, stack string) string {
	data := map[string]interface{}{
		"assertData": assertData,
	}
	if stack != "" {
		data["stack"] = stack
	}
	out, _ := yaml.Marshal(data)
	return string(out)
}
