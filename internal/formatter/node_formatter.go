package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/slurm-rest-cli/srest/pkg/models"
)

// NodeFormatter is a specialized formatter for node data
type NodeFormatter struct{}

// Format formats node data into a readable table
func (f *NodeFormatter) Format(data interface{}) (string, error) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write headers
	fmt.Fprintln(w, "NAME\tSTATE\tCPUS\tSOCKETS\tCORES\tTHREADS\tMEMORY\tFEATURES\tPARTITION\tWEIGHT")

	// Process the data based on its type
	switch nodes := data.(type) {
	case []models.NodeInfo:
		// Handle array of NodeInfo objects
		for _, node := range nodes {
			if err := f.formatNodeInfoRow(w, node); err != nil {
				return "", err
			}
		}
	case models.NodeInfo:
		// Handle single NodeInfo object
		if err := f.formatNodeInfoRow(w, nodes); err != nil {
			return "", err
		}
	case []interface{}:
		// Handle array of nodes
		for _, nodeData := range nodes {
			if err := f.formatNodeRow(w, nodeData); err != nil {
				return "", err
			}
		}
	case map[string]interface{}:
		// Handle single node
		if err := f.formatNodeRow(w, nodes); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported node data type: %T", data)
	}

	// Flush the tabwriter
	if err := w.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatNodeRow formats a single node row
func (f *NodeFormatter) formatNodeRow(w *tabwriter.Writer, nodeData interface{}) error {
	node, ok := nodeData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid node data format: %T", nodeData)
	}

	// Extract node name
	name := safeGetString(node, "name")

	// Extract state
	var state string
	if stateVal, ok := node["state"]; ok {
		switch s := stateVal.(type) {
		case string:
			state = s
		case []interface{}:
			states := make([]string, len(s))
			for i, st := range s {
				states[i] = fmt.Sprintf("%v", st)
			}
			state = strings.Join(states, "+")
		default:
			state = fmt.Sprintf("%v", stateVal)
		}
	}

	// Extract CPU info
	cpus := safeGetInt(node, "cpus")
	sockets := safeGetInt(node, "sockets")
	cores := safeGetInt(node, "cores_per_socket")
	threads := safeGetInt(node, "threads_per_core")

	// Extract memory
	var memory string
	if memVal, ok := node["real_memory"]; ok {
		memInt, ok := memVal.(float64)
		if ok {
			memory = fmt.Sprintf("%.0fMB", memInt)
		} else {
			memory = fmt.Sprintf("%v", memVal)
		}
	}

	// Extract features
	var features string
	if featVal, ok := node["features"]; ok {
		switch f := featVal.(type) {
		case string:
			features = f
		case []interface{}:
			feats := make([]string, len(f))
			for i, feat := range f {
				feats[i] = fmt.Sprintf("%v", feat)
			}
			features = strings.Join(feats, ",")
		default:
			features = fmt.Sprintf("%v", featVal)
		}
	}

	// Extract partition
	partition := safeGetString(node, "partitions")

	// Extract weight
	weight := safeGetInt(node, "weight")

	// Write the row
	fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%d\n",
		name, state, cpus, sockets, cores, threads, memory, features, partition, weight)

	return nil
}

// safeGetString safely extracts a string value from a map
func safeGetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// safeGetInt safely extracts an integer value from a map
func safeGetInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return 0
}

// formatNodeInfoRow formats a single node row from a NodeInfo struct
func (f *NodeFormatter) formatNodeInfoRow(w *tabwriter.Writer, node models.NodeInfo) error {
	// Extract state as string
	var state string
	switch s := node.State.(type) {
	case string:
		state = s
	case []interface{}:
		states := make([]string, len(s))
		for i, st := range s {
			states[i] = fmt.Sprintf("%v", st)
		}
		state = strings.Join(states, "+")
	case nil:
		state = "UNKNOWN"
	default:
		state = fmt.Sprintf("%v", s)
	}

	// Extract features
	var features string
	switch f := node.Features.(type) {
	case string:
		features = f
	case []interface{}:
		feats := make([]string, len(f))
		for i, feat := range f {
			feats[i] = fmt.Sprintf("%v", feat)
		}
		features = strings.Join(feats, ",")
	case nil:
		features = ""
	default:
		features = fmt.Sprintf("%v", f)
	}

	// Format memory
	memory := fmt.Sprintf("%dMB", node.RealMemory)

	// Format partitions
	partition := strings.Join(node.Partitions, ",")

	// Write the row
	fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%d\n",
		node.Name, state, node.CPUs, node.Sockets, node.CoresPerSocket, 
		node.ThreadsPerCore, memory, features, partition, node.Weight)

	return nil
}
