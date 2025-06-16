package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Job submission flags
	jobName      string
	jobAccount   string
	jobPartition string
	jobQos       string
	jobNodes     int
	jobTasks     int
	jobCpusPerTask int
	jobMemPerCpu  string
	jobMemPerNode string
	jobTimeLimit  string
	jobBeginTime  string
	jobStdIn      string
	jobStdOut     string
	jobStdErr     string
	jobHold       bool
	jobRequeue    bool
	jobArray      string
	jobDependency string
	jobConstraint string

	// Job list flags
	jobUser      string
	jobPartitionFilter string
	jobStateFilter     string
	jobFormat          string
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Job-related operations",
	Long:  `Commands for managing Slurm jobs via the REST API.`,
}

// jobSubmitCmd represents the job submit command
var jobSubmitCmd = &cobra.Command{
	Use:   "submit [options] [script]",
	Short: "Submit a batch job",
	Long: `Submit a batch job to the Slurm scheduler.
If a script file is provided, its contents will be used as the job script.
Otherwise, the script will be read from standard input.`,
	Run: func(cmd *cobra.Command, args []string) {
		var scriptContent string
		var err error

		// Read script content from file or stdin
		if len(args) > 0 {
			scriptPath := args[0]
			var scriptBytes []byte
			scriptBytes, err = os.ReadFile(scriptPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading script file: %v\n", err)
				os.Exit(1)
			}
			scriptContent = string(scriptBytes)
		} else {
			// Read from stdin
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
			scriptContent = string(stdinBytes)
		}

		// Execute job submission using the API client
		if err := executeJobSubmit(scriptContent); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// jobListCmd represents the job list command
var jobListCmd = &cobra.Command{
	Use:   "list [options]",
	Short: "List jobs in queue",
	Long:  `List jobs in the Slurm queue with optional filtering.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Execute job list using the API client
		if err := executeJobList(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// jobShowCmd represents the job show command
var jobShowCmd = &cobra.Command{
	Use:   "show <job_id>",
	Short: "Show job details",
	Long:  `Show detailed information about a specific job.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobID := args[0]

		// Execute job show using the API client
		if err := executeJobShow(jobID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// jobCancelCmd represents the job cancel command
var jobCancelCmd = &cobra.Command{
	Use:   "cancel <job_id>",
	Short: "Cancel a job",
	Long:  `Cancel a running or pending job.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobID := args[0]

		// Execute job cancel using the API client
		if err := executeJobCancel(jobID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Helper function to generate curl commands
func generateCurlCommand(method, endpoint string, headers map[string]string, body string) string {
	baseURL := apiURL
	if baseURL == "" {
		baseURL = "https://slurm-api.example.com"
	}

	url := baseURL + endpoint
	headerStr := ""
	for k, v := range headers {
		headerStr += fmt.Sprintf(" -H '%s: %s'", k, v)
	}

	bodyStr := ""
	if body != "" {
		bodyStr = fmt.Sprintf(" -d '%s'", body)
	}

	return fmt.Sprintf("curl -X %s%s '%s'%s", method, headerStr, url, bodyStr)
}

// Helper function to truncate a string
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	rootCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobSubmitCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
	jobCmd.AddCommand(jobCancelCmd)

	// Job submit flags
	jobSubmitCmd.Flags().StringVarP(&jobName, "job-name", "J", "", "Job name")
	jobSubmitCmd.Flags().StringVarP(&jobAccount, "account", "A", "", "Account to charge resources to")
	jobSubmitCmd.Flags().StringVarP(&jobPartition, "partition", "p", "", "Partition to submit job to")
	jobSubmitCmd.Flags().StringVar(&jobQos, "qos", "", "Quality of service")
	jobSubmitCmd.Flags().IntVarP(&jobNodes, "nodes", "N", 1, "Number of nodes required")
	jobSubmitCmd.Flags().IntVarP(&jobTasks, "ntasks", "n", 1, "Number of tasks")
	jobSubmitCmd.Flags().IntVarP(&jobCpusPerTask, "cpus-per-task", "c", 1, "CPUs per task")
	jobSubmitCmd.Flags().StringVar(&jobMemPerCpu, "mem-per-cpu", "", "Memory per CPU")
	jobSubmitCmd.Flags().StringVar(&jobMemPerNode, "mem", "", "Memory per node")
	jobSubmitCmd.Flags().StringVarP(&jobTimeLimit, "time", "t", "", "Time limit")
	jobSubmitCmd.Flags().StringVar(&jobBeginTime, "begin", "", "Begin time")
	jobSubmitCmd.Flags().StringVarP(&jobStdIn, "input", "i", "", "Standard input file")
	jobSubmitCmd.Flags().StringVarP(&jobStdOut, "output", "o", "", "Standard output file")
	jobSubmitCmd.Flags().StringVarP(&jobStdErr, "error", "e", "", "Standard error file")
	jobSubmitCmd.Flags().BoolVarP(&jobHold, "hold", "H", false, "Submit job in held state")
	jobSubmitCmd.Flags().BoolVar(&jobRequeue, "requeue", false, "Requeue job on failure")
	jobSubmitCmd.Flags().StringVarP(&jobArray, "array", "a", "", "Job array indices")
	jobSubmitCmd.Flags().StringVarP(&jobDependency, "dependency", "d", "", "Job dependencies")
	jobSubmitCmd.Flags().StringVarP(&jobConstraint, "constraint", "C", "", "Node constraints")

	// Job list flags
	jobListCmd.Flags().StringVarP(&jobUser, "user", "u", "", "Filter jobs by user")
	jobListCmd.Flags().StringVarP(&jobPartitionFilter, "partition", "p", "", "Filter jobs by partition")
	jobListCmd.Flags().StringVar(&jobStateFilter, "state", "", "Filter jobs by state")
	jobListCmd.Flags().StringVarP(&jobFormat, "format", "o", "", "Output format specifier")
}
