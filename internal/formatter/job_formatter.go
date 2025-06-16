package formatter

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/slurm-rest-cli/srest/pkg/models"
)

// JobFormatter is a specialized formatter for job data
type JobFormatter struct{}

// Format formats job data into a readable table
func (f *JobFormatter) Format(data interface{}) (string, error) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write headers
	fmt.Fprintln(w, "JOBID\tNAME\tUSER\tPARTITION\tSTATE\tTIME\tNODES\tNODELIST")

	// Process the data based on its type
	switch jobs := data.(type) {
	case []models.JobInfo:
		// Handle array of JobInfo objects
		for _, job := range jobs {
			if err := f.formatJobInfoRow(w, job); err != nil {
				return "", err
			}
		}
	case models.JobInfo:
		// Handle single JobInfo object
		if err := f.formatJobInfoRow(w, jobs); err != nil {
			return "", err
		}
	case []interface{}:
		// Handle array of jobs as generic objects
		for _, jobData := range jobs {
			if err := f.formatJobRow(w, jobData); err != nil {
				return "", err
			}
		}
	case map[string]interface{}:
		// Handle single job as generic object
		if err := f.formatJobRow(w, jobs); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported job data type: %T", data)
	}

	// Flush the tabwriter
	if err := w.Flush(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatJobInfoRow formats a single job row from a JobInfo struct
func (f *JobFormatter) formatJobInfoRow(w *tabwriter.Writer, job models.JobInfo) error {
	// Format job state
	state := job.JobState
	if state == "" {
		state = "UNKNOWN"
	}

	// Format job runtime
	timeRunning := "00:00:00"
	if job.StartTime > 0 {
		var endTime int64
		if job.EndTime > 0 {
			endTime = job.EndTime
		} else {
			endTime = time.Now().Unix()
		}
		
		duration := endTime - job.StartTime
		hours := duration / 3600
		minutes := (duration % 3600) / 60
		seconds := duration % 60
		timeRunning = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	// Format nodes
	nodes := job.Nodes
	if nodes == "" {
		nodes = "n/a"
	}

	// Write the row
	fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
		job.JobID, job.Name, job.UserName, job.Partition, state, timeRunning, job.NodeCount, nodes)

	return nil
}

// formatJobRow formats a single job row from a generic map
func (f *JobFormatter) formatJobRow(w *tabwriter.Writer, data interface{}) error {
	jobMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid job data format: %T", data)
	}

	// Extract values with appropriate defaults
	jobID := jobSafeGetInt(jobMap, "job_id")
	name := jobSafeGetString(jobMap, "name")
	userName := jobSafeGetString(jobMap, "user_name")
	partition := jobSafeGetString(jobMap, "partition")
	
	// Format job state
	state := jobSafeGetString(jobMap, "job_state")
	if state == "" {
		state = "UNKNOWN"
	}

	// Format job runtime
	timeRunning := "00:00:00"
	startTime := jobSafeGetInt64(jobMap, "start_time")
	if startTime > 0 {
		var endTime int64
		endTimeVal := jobSafeGetInt64(jobMap, "end_time")
		if endTimeVal > 0 {
			endTime = endTimeVal
		} else {
			endTime = time.Now().Unix()
		}
		
		duration := endTime - startTime
		hours := duration / 3600
		minutes := (duration % 3600) / 60
		seconds := duration % 60
		timeRunning = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	// Format nodes
	nodes := jobSafeGetString(jobMap, "nodes")
	if nodes == "" {
		nodes = "n/a"
	}
	
	nodeCount := jobSafeGetInt(jobMap, "node_count")

	// Write the row
	fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
		jobID, name, userName, partition, state, timeRunning, nodeCount, nodes)

	return nil
}

// jobSafeGetString safely extracts a string value from a map
func jobSafeGetString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// jobSafeGetInt safely extracts an int value from a map
func jobSafeGetInt(data map[string]interface{}, key string) int {
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

// jobSafeGetInt64 safely extracts an int64 value from a map
func jobSafeGetInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}
	return 0
}
