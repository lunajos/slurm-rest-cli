package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Partition list and update flags
	partitionState string
	partitionName  string
	partitionFormat string
	partitionNodes string
	partitionQOS string
	partitionTRES string

	// Partition update flags
	partitionDefaultTime   string
	partitionMaxTime       string // Used for both list filtering and updates
	partitionPriority      int
	partitionAllowGroups   string
	partitionAllowAccounts string
	partitionAllowQOS      string
	partitionDenyAccounts  string
	partitionDenyQOS       string
	partitionOverSubscribe string
	partitionHidden        int
	partitionMaxNodes      int
	partitionMinNodes      int
	partitionDefMemPerCPU  string
	partitionDefMemPerNode string
	partitionMaxMemPerCPU  string
	partitionMaxMemPerNode string
	// partitionNodes is already declared above
	partitionAllocNodes    string
	partitionAlternate     string
	partitionGraceTime     int
	// partitionQOS is already declared above
	partitionDisableRoot   int
	partitionExclusiveUser int
	partitionOverTimeLimit int
	partitionPreemptMode   string
)

// partitionCmd represents the partition command
var partitionCmd = &cobra.Command{
	Use:   "partition",
	Short: "Partition-related operations",
	Long:  `Commands for managing Slurm partitions via the REST API.`,
}

// partitionListCmd represents the partition list command
var partitionListCmd = &cobra.Command{
	Use:   "list [options]",
	Short: "List partition information",
	Long:  `List information about partitions in the Slurm cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Execute partition list using the API client
		if err := executePartitionList(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// partitionShowCmd represents the partition show command
var partitionShowCmd = &cobra.Command{
	Use:   "show <partition_name>",
	Short: "Show partition details",
	Long:  `Show detailed information about a specific partition.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		partitionName := args[0]

		// Execute partition show using the API client
		if err := executePartitionShow(partitionName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// partitionUpdateCmd represents the partition update command
var partitionUpdateCmd = &cobra.Command{
	Use:   "update <partition_name> [options]",
	Short: "Update partition properties",
	Long:  `Update properties of a specific partition.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		partitionName := args[0]

		// Execute partition update using the API client
		if err := executePartitionUpdate(partitionName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(partitionCmd)
	partitionCmd.AddCommand(partitionListCmd)
	partitionCmd.AddCommand(partitionShowCmd)
	partitionCmd.AddCommand(partitionUpdateCmd)

	// Partition list flags
	partitionListCmd.Flags().StringVar(&partitionState, "state", "", "Filter partitions by state")
	partitionListCmd.Flags().StringVar(&partitionName, "name", "", "Filter partitions by name")
	partitionListCmd.Flags().StringVarP(&partitionFormat, "format", "o", "", "Output format specifier")
	partitionListCmd.Flags().StringVar(&partitionNodes, "nodes", "", "Filter partitions by node list")
	partitionListCmd.Flags().StringVar(&partitionQOS, "qos", "", "Filter partitions by QOS")
	partitionListCmd.Flags().StringVar(&partitionTRES, "tres", "", "Filter partitions by TRES")
	partitionListCmd.Flags().StringVar(&partitionMaxTime, "max-time", "", "Filter partitions by max time limit")

	// Partition update flags
	partitionUpdateCmd.Flags().StringVar(&partitionState, "state", "", "Set partition state (e.g., UP, DOWN, DRAIN, INACTIVE)")
	partitionUpdateCmd.Flags().StringVar(&partitionDefaultTime, "default-time", "", "Default time limit (e.g., '1:00:00' for 1 hour)")
	partitionUpdateCmd.Flags().StringVar(&partitionMaxTime, "max-time", "", "Maximum time limit (e.g., '24:00:00' for 24 hours)")
	partitionUpdateCmd.Flags().IntVar(&partitionPriority, "priority", 0, "Partition priority")
	partitionUpdateCmd.Flags().StringVar(&partitionAllowGroups, "allow-groups", "", "Allowed groups (comma-separated)")
	partitionUpdateCmd.Flags().StringVar(&partitionAllowAccounts, "allow-accounts", "", "Allowed accounts (comma-separated)")
	partitionUpdateCmd.Flags().StringVar(&partitionAllowQOS, "allow-qos", "", "Allowed QOS (comma-separated)")
	partitionUpdateCmd.Flags().StringVar(&partitionDenyAccounts, "deny-accounts", "", "Denied accounts (comma-separated)")
	partitionUpdateCmd.Flags().StringVar(&partitionDenyQOS, "deny-qos", "", "Denied QOS (comma-separated)")
	partitionUpdateCmd.Flags().StringVar(&partitionOverSubscribe, "oversubscribe", "", "Oversubscribe setting (YES, NO, EXCLUSIVE, FORCE)")
	partitionUpdateCmd.Flags().IntVar(&partitionHidden, "hidden", -1, "Hidden flag (1=hidden, 0=visible, -1=unchanged)")
	partitionUpdateCmd.Flags().IntVar(&partitionMaxNodes, "max-nodes", 0, "Maximum nodes per job (0=unlimited)")
	partitionUpdateCmd.Flags().IntVar(&partitionMinNodes, "min-nodes", 0, "Minimum nodes per job")
	partitionUpdateCmd.Flags().StringVar(&partitionDefMemPerCPU, "def-mem-per-cpu", "", "Default memory per CPU (e.g., '1G')")
	partitionUpdateCmd.Flags().StringVar(&partitionDefMemPerNode, "def-mem-per-node", "", "Default memory per node (e.g., '4G')")
	partitionUpdateCmd.Flags().StringVar(&partitionMaxMemPerCPU, "max-mem-per-cpu", "", "Maximum memory per CPU (e.g., '2G')")
	partitionUpdateCmd.Flags().StringVar(&partitionMaxMemPerNode, "max-mem-per-node", "", "Maximum memory per node (e.g., '64G')")
	partitionUpdateCmd.Flags().StringVar(&partitionNodes, "nodes", "", "Node list for partition")
	partitionUpdateCmd.Flags().StringVar(&partitionAllocNodes, "alloc-nodes", "", "Allowed allocation nodes")
	partitionUpdateCmd.Flags().StringVar(&partitionAlternate, "alternate", "", "Alternate partition")
	partitionUpdateCmd.Flags().IntVar(&partitionGraceTime, "grace-time", 0, "Grace time for jobs")
	partitionUpdateCmd.Flags().StringVar(&partitionQOS, "qos", "", "Quality of Service")
	partitionUpdateCmd.Flags().IntVar(&partitionDisableRoot, "disable-root-jobs", -1, "Disable root jobs (1=disable, 0=allow, -1=unchanged)")
	partitionUpdateCmd.Flags().IntVar(&partitionExclusiveUser, "exclusive-user", -1, "Exclusive user (1=exclusive, 0=shared, -1=unchanged)")
	partitionUpdateCmd.Flags().IntVar(&partitionOverTimeLimit, "over-time-limit", -1, "Over time limit")
	partitionUpdateCmd.Flags().StringVar(&partitionPreemptMode, "preempt-mode", "", "Preemption mode")
}
