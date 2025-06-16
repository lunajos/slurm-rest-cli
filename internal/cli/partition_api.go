package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/slurm-rest-cli/srest/internal/api"
	"github.com/slurm-rest-cli/srest/internal/formatter"
	"github.com/slurm-rest-cli/srest/pkg/config"
	"github.com/slurm-rest-cli/srest/pkg/models"
	"github.com/spf13/viper"
)

// createPartitionClient creates a new partition client using the current configuration.
func createPartitionClient() (*api.PartitionClient, error) {
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
	return api.NewPartitionClient(client), nil
}

// executePartitionList executes the partition list command with the API client.
func executePartitionList() error {
	// Create partition client
	partitionClient, err := createPartitionClient()
	if err != nil {
		return err
	}

	// Build filters
	filters := make(map[string]string)
	if partitionState != "" {
		filters["state"] = partitionState
	}
	if partitionName != "" {
		filters["name"] = partitionName
	}
	if partitionNodes != "" {
		filters["nodes"] = partitionNodes
	}
	if partitionQOS != "" {
		filters["qos"] = partitionQOS
	}
	if partitionTRES != "" {
		filters["tres"] = partitionTRES
	}
	if partitionMaxTime != "" {
		filters["max_time"] = partitionMaxTime
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := partitionClient.GenerateListCurl(filters)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// List the partitions
	response, err := partitionClient.List(filters)
	if err != nil {
		return fmt.Errorf("error listing partitions: %w", err)
	}

	// Check for errors in the response but don't fail if they're just warnings
	if len(response.Errors) > 0 {
		// Print warnings but don't fail if we have partition data
		fmt.Fprintf(os.Stderr, "Warning: %v\n", response.Errors)
	}
	
	// Check if we have any partition data
	if response.Partitions == nil || len(response.Partitions) == 0 {
		return fmt.Errorf("no partition data returned")
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewTypedFormatter(fmtType, formatter.PartitionDataType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Partitions)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executePartitionShow executes the partition show command with the API client.
func executePartitionShow(partitionName string) error {
	// Create partition client
	partitionClient, err := createPartitionClient()
	if err != nil {
		return err
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := partitionClient.GenerateGetCurl(partitionName)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Get the partition
	response, err := partitionClient.Get(partitionName)
	if err != nil {
		return fmt.Errorf("error getting partition: %w", err)
	}

	// Check for errors in the response, but don't fail if they're just warnings or mock data
	if len(response.Errors) > 0 {
		// Check if this is mock data (we don't want to fail in this case)
		isMockData := false
		for _, err := range response.Errors {
			if strings.Contains(err, "mock data") {
				isMockData = true
				break
			}
		}
		
		if !isMockData {
			return fmt.Errorf("partition show failed: %s", strings.Join(response.Errors, ", "))
		} else {
			// Just print a warning for mock data
			fmt.Fprintf(os.Stderr, "Note: Using mock partition data due to authentication issues\n")
		}
	}

	// Format the output
	outputFormat := viper.GetString("preferences.output_format")
	if outputFormat == "" {
		outputFormat = "table"
	}

	fmtType := formatter.FormatType(outputFormat)
	formatterObj, err := formatter.NewTypedFormatter(fmtType, formatter.PartitionDataType)
	if err != nil {
		return fmt.Errorf("error creating formatter: %w", err)
	}

	output, err := formatterObj.Format(response.Partition)
	if err != nil {
		return fmt.Errorf("error formatting response: %w", err)
	}

	// Print the output
	fmt.Println(output)
	return nil
}

// executePartitionUpdate executes the partition update command with the API client.
func executePartitionUpdate(partitionName string) error {
	// Create partition client
	partitionClient, err := createPartitionClient()
	if err != nil {
		return err
	}

	// Create partition update request
	request := &models.PartitionUpdateRequest{
		State:           partitionState,
		DefaultTime:     partitionDefaultTime,
		MaxTime:         partitionMaxTime,
		Priority:        partitionPriority,
		AllowGroups:     partitionAllowGroups,
		AllowAccounts:   partitionAllowAccounts,
		AllowQOS:        partitionAllowQOS,
		DenyAccounts:    partitionDenyAccounts,
		DenyQOS:         partitionDenyQOS,
		OverSubscribe:   partitionOverSubscribe,
		Hidden:          partitionHidden,
		MaxNodes:        partitionMaxNodes,
		MinNodes:        partitionMinNodes,
		Nodes:           partitionNodes,
		AllocNodes:      partitionAllocNodes,
		Alternate:       partitionAlternate,
		GraceTime:       partitionGraceTime,
		QOS:             partitionQOS,
		DisableRootJobs: partitionDisableRoot,
		ExclusiveUser:   partitionExclusiveUser,
		OverTimeLimit:   partitionOverTimeLimit,
		PreemptMode:     partitionPreemptMode,
	}

	// If DefMemPerCPU flag is set
	if partitionDefMemPerCPU != "" {
		// Parse memory value (e.g., "1G", "100M") - simplified for now
		var memValue int64
		fmt.Sscanf(partitionDefMemPerCPU, "%d", &memValue)
		request.DefMemPerCPU = memValue
	}

	// If DefMemPerNode flag is set
	if partitionDefMemPerNode != "" {
		// Parse memory value (e.g., "1G", "100M") - simplified for now
		var memValue int64
		fmt.Sscanf(partitionDefMemPerNode, "%d", &memValue)
		request.DefMemPerNode = memValue
	}

	// If MaxMemPerCPU flag is set
	if partitionMaxMemPerCPU != "" {
		// Parse memory value (e.g., "1G", "100M") - simplified for now
		var memValue int64
		fmt.Sscanf(partitionMaxMemPerCPU, "%d", &memValue)
		request.MaxMemPerCPU = memValue
	}

	// If MaxMemPerNode flag is set
	if partitionMaxMemPerNode != "" {
		// Parse memory value (e.g., "1G", "100M") - simplified for now
		var memValue int64
		fmt.Sscanf(partitionMaxMemPerNode, "%d", &memValue)
		request.MaxMemPerNode = memValue
	}

	// If curl flag is set, just print the curl command
	if showCurl {
		curlCmd, err := partitionClient.GenerateUpdateCurl(partitionName, request)
		if err != nil {
			return fmt.Errorf("error generating curl command: %w", err)
		}
		fmt.Println("Equivalent curl command:")
		fmt.Println(curlCmd)
		return nil
	}

	// Update the partition
	if err := partitionClient.Update(partitionName, request); err != nil {
		return fmt.Errorf("error updating partition: %w", err)
	}

	fmt.Printf("Partition %s updated successfully\n", partitionName)
	return nil
}
