package formatter

import (
	"fmt"
)

// Formatter is an interface for formatting data in different output formats.
type Formatter interface {
	Format(data interface{}) (string, error)
}

// FormatType represents the type of formatter to use.
type FormatType string

const (
	// JSONFormat represents JSON output format.
	JSONFormat FormatType = "json"
	// YAMLFormat represents YAML output format.
	YAMLFormat FormatType = "yaml"
	// TableFormat represents table output format.
	TableFormat FormatType = "table"
)

// DataType represents the type of data to format
type DataType string

const (
	// NodeDataType represents node data
	NodeDataType DataType = "node"
	// PartitionDataType represents partition data
	PartitionDataType DataType = "partition"
	// JobDataType represents job data
	JobDataType DataType = "job"
	// GenericDataType represents generic data
	GenericDataType DataType = "generic"
)

// NewFormatter creates a new formatter of the specified type.
func NewFormatter(formatType FormatType) (Formatter, error) {
	switch formatType {
	case JSONFormat:
		return &JSONFormatter{}, nil
	case YAMLFormat:
		return &YAMLFormatter{}, nil
	case TableFormat:
		return &TableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format type: %s", formatType)
	}
}

// NewTypedFormatter creates a new formatter for a specific data type
func NewTypedFormatter(formatType FormatType, dataType DataType) (Formatter, error) {
	// For JSON and YAML, use the standard formatters
	if formatType == JSONFormat || formatType == YAMLFormat {
		return NewFormatter(formatType)
	}
	
	// For table format, use specialized formatters based on data type
	if formatType == TableFormat {
		switch dataType {
		case NodeDataType:
			return &NodeFormatter{}, nil
		case PartitionDataType:
			return &PartitionFormatter{}, nil
		case JobDataType:
			return &JobFormatter{}, nil
		case GenericDataType:
			// Use the standard table formatter for generic data
			return &TableFormatter{}, nil
		default:
			return &TableFormatter{}, nil
		}
	}
	
	return nil, fmt.Errorf("unsupported format type: %s", formatType)
}
