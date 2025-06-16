package formatter

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"
)

// TableFormatter formats data as a table.
type TableFormatter struct {
	// Headers are the column headers for the table.
	Headers []string
	// Fields are the field names to extract from the data.
	Fields []string
}

// Format formats the data as a table.
func (f *TableFormatter) Format(data interface{}) (string, error) {
	// Create a buffer for the table output
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write headers if provided
	if len(f.Headers) > 0 {
		fmt.Fprintln(w, strings.Join(f.Headers, "\t"))
	}

	// Process the data based on its type
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		// Handle slice/array of structs or maps
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if err := f.formatRow(w, item); err != nil {
				return "", err
			}
		}
	case reflect.Map, reflect.Struct:
		// Handle single struct or map
		if err := f.formatRow(w, data); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported data type for table formatting: %T", data)
	}

	// Flush the tabwriter
	if err := w.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatRow formats a single row of the table.
func (f *TableFormatter) formatRow(w *tabwriter.Writer, data interface{}) error {
	// If no fields are specified, try to infer them
	fields := f.Fields
	if len(fields) == 0 {
		var err error
		fields, err = inferFields(data)
		if err != nil {
			return err
		}
	}

	// Extract values for each field
	values := make([]string, len(fields))
	for i, field := range fields {
		value, err := extractValue(data, field)
		if err != nil {
			values[i] = "<error>"
		} else {
			values[i] = fmt.Sprintf("%v", value)
		}
	}

	// Write the row
	fmt.Fprintln(w, strings.Join(values, "\t"))
	return nil
}

// inferFields attempts to infer field names from the data.
func inferFields(data interface{}) ([]string, error) {
	v := reflect.ValueOf(data)
	
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		// For structs, use the field names
		t := v.Type()
		fields := make([]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fields[i] = t.Field(i).Name
		}
		return fields, nil
	case reflect.Map:
		// For maps, use the keys
		keys := v.MapKeys()
		fields := make([]string, len(keys))
		for i, key := range keys {
			fields[i] = fmt.Sprintf("%v", key.Interface())
		}
		return fields, nil
	default:
		return nil, fmt.Errorf("cannot infer fields from type %T", data)
	}
}

// extractValue extracts a value from the data using the field name.
func extractValue(data interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(data)
	
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		// For structs, get the field by name
		f := v.FieldByName(field)
		if !f.IsValid() {
			return nil, fmt.Errorf("field %s not found", field)
		}
		return f.Interface(), nil
	case reflect.Map:
		// For maps, get the value by key
		key := reflect.ValueOf(field)
		value := v.MapIndex(key)
		if !value.IsValid() {
			return nil, fmt.Errorf("key %s not found", field)
		}
		return value.Interface(), nil
	default:
		return nil, fmt.Errorf("cannot extract value from type %T", data)
	}
}
