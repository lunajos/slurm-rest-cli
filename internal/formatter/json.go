package formatter

import (
	"encoding/json"
)

// JSONFormatter formats data as JSON.
type JSONFormatter struct {
	// PrettyPrint controls whether to format the JSON with indentation.
	PrettyPrint bool
}

// Format formats the data as JSON.
func (f *JSONFormatter) Format(data interface{}) (string, error) {
	var (
		bytes []byte
		err   error
	)

	if f.PrettyPrint {
		bytes, err = json.MarshalIndent(data, "", "  ")
	} else {
		bytes, err = json.Marshal(data)
	}

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
