package formatter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/slurm-rest-cli/srest/pkg/models"
)

// PartitionFormatter is a specialized formatter for partition data
type PartitionFormatter struct{}

// Format formats partition data into a readable table
func (f *PartitionFormatter) Format(data interface{}) (string, error) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write headers
	fmt.Fprintln(w, "NAME\tSTATE\tDEFAULT\tNODES\tTIME_LIMIT\tMAX_CPUS\tMIN_NODES\tMAX_NODES\tPRIORITY")

	// Process the data based on its type
	switch partitions := data.(type) {
	case []models.PartitionInfo:
		// Handle array of PartitionInfo objects
		for _, partition := range partitions {
			if err := f.formatPartitionInfoRow(w, partition); err != nil {
				return "", err
			}
		}
	case models.PartitionInfo:
		// Handle single PartitionInfo object
		if err := f.formatPartitionInfoRow(w, partitions); err != nil {
			return "", err
		}
	case []interface{}:
		// Handle array of partitions as generic objects
		for _, partitionData := range partitions {
			if err := f.formatPartitionRow(w, partitionData); err != nil {
				return "", err
			}
		}
	case map[string]interface{}:
		// Handle single partition as generic object
		if err := f.formatPartitionRow(w, partitions); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported partition data type: %T", data)
	}

	// Flush the tabwriter
	if err := w.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatPartitionInfoRow formats a single partition row from a PartitionInfo struct
func (f *PartitionFormatter) formatPartitionInfoRow(w *tabwriter.Writer, partition models.PartitionInfo) error {
	// Format state
	state := partition.State
	if state == "" {
		state = "UP"
	}

	// Format default flag
	defaultFlag := "no"
	if partition.Flags != nil {
		for _, flag := range partition.Flags {
			if flag == "Default" {
				defaultFlag = "yes"
				break
			}
		}
	}

	// Format time limit
	timeLimit := "infinite"
	if partition.MaxTime != "" && partition.MaxTime != "UNLIMITED" {
		// Try to parse the time as minutes
		minutes, err := strconv.Atoi(partition.MaxTime)
		if err == nil && minutes > 0 {
			hours := minutes / 60
			mins := minutes % 60
			timeLimit = fmt.Sprintf("%02d:%02d:00", hours, mins)
		} else {
			// If not a simple number, use as is
			timeLimit = partition.MaxTime
		}
	}

	// Format nodes
	nodes := partition.NodeList
	if nodes == "" {
		nodes = "n/a"
	}

	// Write the row
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\n",
		partition.Name, state, defaultFlag, nodes, timeLimit,
		partition.MaxCPUsPerNode, partition.MinNodes, partition.MaxNodes, partition.Priority)

	return nil
}

// formatPartitionRow formats a single partition row from a generic map
func (f *PartitionFormatter) formatPartitionRow(w *tabwriter.Writer, data interface{}) error {
	partitionMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid partition data format: %T", data)
	}

	// Extract values with appropriate defaults
	name := partitionSafeGetString(partitionMap, "name")
	state := partitionSafeGetString(partitionMap, "state")
	if state == "" {
		state = "UP"
	}

	// Format default flag
	defaultFlag := "no"
	if flags, ok := partitionMap["flags"]; ok {
		if flagsArr, ok := flags.([]interface{}); ok {
			for _, flag := range flagsArr {
				if flagStr, ok := flag.(string); ok && flagStr == "Default" {
					defaultFlag = "yes"
					break
				}
			}
		} else if flagStr, ok := flags.(string); ok && strings.Contains(flagStr, "Default") {
			defaultFlag = "yes"
		}
	}

	// Format time limit
	timeLimit := "infinite"
	if maxTime, ok := partitionMap["max_time"]; ok {
		if mt, ok := maxTime.(float64); ok && mt > 0 {
			hours := int(mt) / 60
			mins := int(mt) % 60
			timeLimit = fmt.Sprintf("%02d:%02d:00", hours, mins)
		}
	}

	// Format nodes
	nodes := "n/a"
	if nodeData, ok := partitionMap["nodes"]; ok {
		switch n := nodeData.(type) {
		case string:
			nodes = n
		case []interface{}:
			nodeList := make([]string, len(n))
			for i, node := range n {
				nodeList[i] = fmt.Sprintf("%v", node)
			}
			nodes = strings.Join(nodeList, ",")
		}
	}

	// Write the row
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\n",
		name, state, defaultFlag, nodes, timeLimit,
		partitionSafeGetInt(partitionMap, "max_cpus_per_node"),
		partitionSafeGetInt(partitionMap, "min_nodes"),
		partitionSafeGetInt(partitionMap, "max_nodes"),
		partitionSafeGetInt(partitionMap, "priority"))

	return nil
}

// partitionSafeGetString safely extracts a string value from a map
func partitionSafeGetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// partitionSafeGetInt safely extracts an int value from a map
func partitionSafeGetInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}
