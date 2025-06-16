package cli

import (
	"fmt"
	"strings"

	"github.com/slurm-rest-cli/srest/internal/api"
	"github.com/slurm-rest-cli/srest/internal/formatter"
	"github.com/slurm-rest-cli/srest/pkg/config"
	"github.com/slurm-rest-cli/srest/pkg/models"
	"github.com/spf13/viper"
)

// createNodeClient creates a new node client using the current configuration.
func createNodeClient() (*api.NodeClient, error) {
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
	return api.NewNodeClient(client), nil
}

// executeNodeList executes the node list command with the API client.
func executeNodeList() error {
	// Create node client
	nodeClient, err := createNodeClient()
	if err != nil {
		return err
	}

	// Build filters
	filters := make(map[string]string)
	if nodeState != "" {
		filters["state"] = nodeState
	}
	if nodePartitionFilter != "" {
		filters["partition"] = nodePartitionFilter
	}
	if nodeName != "" {
		filters["name"] = nodeName
	}
	if nodeAddress != "" {
		filters["address"] = nodeAddress
	}
	if nodeHostname != "" {
		filters["hostname"] = nodeHostname
	}
	if nodeFeatureFilter != "" {
		filters["features"] = nodeFeatureFilter
	}
	if nodeGres != "" {
		filters["gres"] = nodeGres
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := nodeClient.GenerateListCurl(filters)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// List the nodes
	response, err := nodeClient.List(filters)
	if err != nil {
		return fmt.Errorf("error listing nodes: %w", err)
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		return fmt.Errorf("node list failed: %s", strings.Join(response.Errors, ", "))
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewTypedFormatter(fmtType, formatter.NodeDataType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Nodes)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executeNodeShow executes the node show command with the API client.
func executeNodeShow(nodeName string) error {
	// Create node client
	nodeClient, err := createNodeClient()
	if err != nil {
		return err
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := nodeClient.GenerateGetCurl(nodeName)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Get the node
	response, err := nodeClient.Get(nodeName)
	if err != nil {
		return fmt.Errorf("error getting node: %w", err)
	}

	// Check for errors in the response
	if len(response.Errors) > 0 {
		return fmt.Errorf("node show failed: %s", strings.Join(response.Errors, ", "))
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewTypedFormatter(fmtType, formatter.NodeDataType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Node)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executeNodeUpdate executes the node update command with the API client.
func executeNodeUpdate(nodeName string) error {
	// Create node client
	nodeClient, err := createNodeClient()
	if err != nil {
		return err
	}

	// Create node update request
	request := &models.NodeUpdateRequest{
		State:    nodeState,
		Reason:   nodeReason,
		Weight:   nodeWeight,
		Features: nodeFeatures,
		Comment:  nodeComment,
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := nodeClient.GenerateUpdateCurl(nodeName, request)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Update the node
	if err := nodeClient.Update(nodeName, request); err != nil {
		return fmt.Errorf("error updating node: %w", err)
	}

	fmt.Printf("Node %s updated successfully\n", nodeName)
	return nil
}
