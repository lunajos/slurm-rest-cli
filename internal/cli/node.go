package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Node list flags
	nodeState          string
	nodePartitionFilter string
	nodeFormat         string
	nodeName           string
	nodeAddress        string
	nodeHostname       string
	nodeFeatureFilter  string
	nodeGres           string
	
	// Node update flags
	nodeReason   string
	nodeWeight   int
	nodeFeatures string
	nodeComment  string
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node-related operations",
	Long:  `Commands for managing Slurm nodes via the REST API.`,
}

// nodeListCmd represents the node list command
var nodeListCmd = &cobra.Command{
	Use:   "list [options]",
	Short: "List node information",
	Long:  `List information about nodes in the Slurm cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Execute node list using the API client
		if err := executeNodeList(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// nodeShowCmd represents the node show command
var nodeShowCmd = &cobra.Command{
	Use:   "show <node_name>",
	Short: "Show node details",
	Long:  `Show detailed information about a specific node.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeName := args[0]

		// Execute node show using the API client
		if err := executeNodeShow(nodeName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// nodeUpdateCmd represents the node update command
var nodeUpdateCmd = &cobra.Command{
	Use:   "update <node_name> [options]",
	Short: "Update node properties",
	Long:  `Update properties of a specific node.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeName := args[0]

		// Execute node update using the API client
		if err := executeNodeUpdate(nodeName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeListCmd)
	nodeCmd.AddCommand(nodeShowCmd)
	nodeCmd.AddCommand(nodeUpdateCmd)

	// Node list flags
	nodeListCmd.Flags().StringVar(&nodeState, "state", "", "Filter nodes by state")
	nodeListCmd.Flags().StringVarP(&nodePartitionFilter, "partition", "p", "", "Filter nodes by partition")
	nodeListCmd.Flags().StringVarP(&nodeFormat, "format", "o", "", "Output format specifier")
	nodeListCmd.Flags().StringVar(&nodeName, "name", "", "Filter nodes by name (supports wildcards)")
	nodeListCmd.Flags().StringVar(&nodeAddress, "address", "", "Filter nodes by address")
	nodeListCmd.Flags().StringVar(&nodeHostname, "hostname", "", "Filter nodes by hostname")
	nodeListCmd.Flags().StringVar(&nodeFeatureFilter, "features", "", "Filter nodes by features")
	nodeListCmd.Flags().StringVar(&nodeGres, "gres", "", "Filter nodes by generic resources")

	// Node update flags
	nodeUpdateCmd.Flags().StringVar(&nodeState, "state", "", "Set node state (e.g., DOWN, DRAIN, RESUME)")
	nodeUpdateCmd.Flags().StringVar(&nodeReason, "reason", "", "Reason for state change")
	nodeUpdateCmd.Flags().IntVar(&nodeWeight, "weight", 0, "Node weight for scheduling")
	nodeUpdateCmd.Flags().StringVar(&nodeFeatures, "features", "", "Node features (comma-separated)")
	nodeUpdateCmd.Flags().StringVar(&nodeComment, "comment", "", "Comment for the node")
}
