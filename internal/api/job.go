package api

import (
	"fmt"
	"net/http"

	"github.com/slurm-rest-cli/srest/pkg/models"
)

// JobClient provides methods for interacting with the Slurm job API.
type JobClient struct {
	client *Client
}

// NewJobClient creates a new job client.
func NewJobClient(client *Client) *JobClient {
	return &JobClient{client: client}
}

// Submit submits a job to the Slurm cluster.
func (c *JobClient) Submit(request *models.JobSubmitRequest) (*models.JobSubmitResponse, error) {
	var response models.JobSubmitResponse
	err := c.client.Request(http.MethodPost, "/slurm/v0.0.42/job/submit", request, &response)
	if err != nil {
		return nil, fmt.Errorf("error submitting job: %w", err)
	}
	return &response, nil
}

// List returns a list of jobs matching the specified filters.
func (c *JobClient) List(filters map[string]string) (*models.JobsResponse, error) {
	// Build query parameters
	endpoint := "/slurm/v0.0.42/jobs"
	
	// Add query parameters for filters
	if len(filters) > 0 {
		endpoint += "?"
		for k, v := range filters {
			endpoint += k + "=" + v + "&"
		}
		// Remove trailing &
		endpoint = endpoint[:len(endpoint)-1]
	}
	
	var response models.JobsResponse
	err := c.client.Request(http.MethodGet, endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error listing jobs: %w", err)
	}
	return &response, nil
}

// Get returns details about a specific job.
func (c *JobClient) Get(jobID int) (*models.JobResponse, error) {
	var response models.JobResponse
	endpoint := fmt.Sprintf("/slurm/v0.0.42/job/%d", jobID)
	err := c.client.Request(http.MethodGet, endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting job %d: %w", jobID, err)
	}
	return &response, nil
}

// Cancel cancels a job.
func (c *JobClient) Cancel(jobID int, signal string) error {
	endpoint := fmt.Sprintf("/slurm/v0.0.42/job/%d", jobID)
	
	// Build request body with signal if provided
	var requestBody map[string]string
	if signal != "" {
		requestBody = map[string]string{"signal": signal}
	}
	
	err := c.client.Request(http.MethodDelete, endpoint, requestBody, nil)
	if err != nil {
		return fmt.Errorf("error canceling job %d: %w", jobID, err)
	}
	return nil
}

// GenerateSubmitCurl generates a curl command for submitting a job.
func (c *JobClient) GenerateSubmitCurl(request *models.JobSubmitRequest) (string, error) {
	return c.client.GenerateCurlCommand(http.MethodPost, "/slurm/v0.0.42/job/submit", request)
}

// GenerateListCurl generates a curl command for listing jobs.
func (c *JobClient) GenerateListCurl(filters map[string]string) (string, error) {
	// Build query parameters
	endpoint := "/slurm/v0.0.42/jobs"
	
	// Add query parameters for filters
	if len(filters) > 0 {
		endpoint += "?"
		for k, v := range filters {
			endpoint += k + "=" + v + "&"
		}
		// Remove trailing &
		endpoint = endpoint[:len(endpoint)-1]
	}
	
	return c.client.GenerateCurlCommand(http.MethodGet, endpoint, nil)
}

// GenerateGetCurl generates a curl command for getting job details.
func (c *JobClient) GenerateGetCurl(jobID int) (string, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.42/job/%d", jobID)
	return c.client.GenerateCurlCommand(http.MethodGet, endpoint, nil)
}

// GenerateCancelCurl generates a curl command for canceling a job.
func (c *JobClient) GenerateCancelCurl(jobID int, signal string) (string, error) {
	endpoint := fmt.Sprintf("/slurm/v0.0.42/job/%d", jobID)
	
	// Build request body with signal if provided
	var requestBody map[string]string
	if signal != "" {
		requestBody = map[string]string{"signal": signal}
	}
	
	return c.client.GenerateCurlCommand(http.MethodDelete, endpoint, requestBody)
}
