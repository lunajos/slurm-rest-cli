package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/slurm-rest-cli/srest/internal/api"
	"github.com/slurm-rest-cli/srest/internal/formatter"
	"github.com/slurm-rest-cli/srest/pkg/config"
	"github.com/slurm-rest-cli/srest/pkg/models"
	"github.com/spf13/viper"
)

// createJobClient creates a new job client using the current configuration.
func createJobClient() (*api.JobClient, error) {
	// Get configuration
	cfg, err := config.LoadFromViper()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	// Get API URL
	apiURL := cfg.GetAPIURL()
	if apiURL == "" {
		return nil, fmt.Errorf("API URL not configured. Use --url flag or set in config")
	}

	// Get authentication info
	token, jwt, username := cfg.GetAuthInfo()
	auth := api.NewAuthenticator(username, token, jwt)

	// Create API client
	client := api.NewClient(apiURL, auth)
	return api.NewJobClient(client), nil
}

// executeJobSubmit executes the job submit command with the API client.
func executeJobSubmit(scriptContent string) error {
	// Create job client
	jobClient, err := createJobClient()
	if err != nil {
		return err
	}

	// Create job submission request
	jobDescr := &models.JobDescr{
		Name:        jobName,
		Account:     jobAccount,
		Partition:   jobPartition,
		QOS:         jobQos,
		Nodes:       jobNodes,
		Tasks:       jobTasks,
		CPUsPerTask: jobCpusPerTask,
		MemPerCPU:   jobMemPerCpu,
		MemPerNode:  jobMemPerNode,
		TimeLimit:   jobTimeLimit,
		BeginTime:   jobBeginTime,
		StdIn:       jobStdIn,
		StdOut:      jobStdOut,
		StdErr:      jobStdErr,
		Hold:        jobHold,
		Requeue:     jobRequeue,
		ArrayInx:    jobArray,
		Dependency:  jobDependency,
		Constraints: jobConstraint,
	}

	request := &models.JobSubmitRequest{
		Script: scriptContent,
		Job:    jobDescr,
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := jobClient.GenerateSubmitCurl(request)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Submit the job
	response, err := jobClient.Submit(request)
	if err != nil {
		return fmt.Errorf("error submitting job: %w", err)
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		return fmt.Errorf("job submission failed: %s", strings.Join(response.Errors, ", "))
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewFormatter(fmtType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executeJobList executes the job list command with the API client.
func executeJobList() error {
	// Create job client
	jobClient, err := createJobClient()
	if err != nil {
		return err
	}

	// Build filters
	filters := make(map[string]string)
	if jobUser != "" {
		filters["user"] = jobUser
	}
	if jobPartitionFilter != "" {
		filters["partition"] = jobPartitionFilter
	}
	if jobStateFilter != "" {
		filters["state"] = jobStateFilter
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := jobClient.GenerateListCurl(filters)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// List the jobs
	response, err := jobClient.List(filters)
	if err != nil {
		return fmt.Errorf("error listing jobs: %w", err)
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		return fmt.Errorf("job list failed: %s", strings.Join(response.Errors, ", "))
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewFormatter(fmtType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Jobs)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executeJobShow executes the job show command with the API client.
func executeJobShow(jobIDStr string) error {
	// Create job client
	jobClient, err := createJobClient()
	if err != nil {
		return err
	}

	// Parse job ID
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		return fmt.Errorf("invalid job ID: %s", jobIDStr)
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := jobClient.GenerateGetCurl(jobID)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Get the job
	response, err := jobClient.Get(jobID)
	if err != nil {
		return fmt.Errorf("error getting job: %w", err)
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		return fmt.Errorf("job show failed: %s", strings.Join(response.Errors, ", "))
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewFormatter(fmtType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Job)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executeJobCancel executes the job cancel command with the API client.
func executeJobCancel(jobIDStr string) error {
	// Create job client
	jobClient, err := createJobClient()
	if err != nil {
		return err
	}

	// Parse job ID
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		return fmt.Errorf("invalid job ID: %s", jobIDStr)
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := jobClient.GenerateCancelCurl(jobID, "")
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Cancel the job
	if err := jobClient.Cancel(jobID, ""); err != nil {
		return fmt.Errorf("error canceling job: %w", err)
	}

	fmt.Printf("Job %d canceled successfully\n", jobID)
	return nil
}
