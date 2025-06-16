package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Job submission flags
	jobName              string
	jobAccount           string
	jobPartition         string
	jobQos               string
	jobComment           string
	jobNodes             string
	jobTasks             int
	jobCpusPerTask       int
	jobMemPerCpu         string
	jobMemPerNode        string
	jobTimeLimit         string
	jobBeginTime         string
	jobDeadline          string
	jobStdIn             string
	jobStdOut            string
	jobStdErr            string
	jobHold              bool
	jobRequeue           bool
	jobKillOnNodeFail    bool
	jobArray             string
	jobDependency        string
	jobMailType          string
	jobMailUser          string
	jobNice              int
	jobConstraint        string
	jobX11               bool
	jobGres              string
	jobTresPerJob        string
	jobTresPerNode       string
	jobTresPerSocket     string
	jobTresPerTask       string
	jobCpusPerTres       string
	jobMemPerTres        string
	jobLicenses          string
	jobClusters          []string
	jobReservation       string
	jobPriority          int
	jobChdir             string
	jobWorkdir           string
	jobWckey             string
	jobExclusive         bool
	jobShared            bool
	jobOversubscribe     bool
	jobContiguous        bool
	jobCoreSpec          int
	jobThreadSpec        int
	jobMinCpus           int
	jobMinNodes          int
	jobMaxNodes          int
	jobSocketsPerNode    int
	jobCoresPerSocket    int
	jobThreadsPerCore    int
	jobNtasksPerNode     int
	jobNtasksPerSocket   int
	jobNtasksPerCore     int
	jobNtasksPerTres     string
	jobEnvironment       []string
	jobBurstBuffer       string
	jobDelayBoot         int
	jobNetwork           string

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
	// Add commands
	rootCmd.AddCommand(jobCmd)
	
	// Add verbose flag for debugging
	jobCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output for debugging")
	jobCmd.AddCommand(jobSubmitCmd)
	jobCmd.AddCommand(jobListCmd)
	jobCmd.AddCommand(jobShowCmd)
	jobCmd.AddCommand(jobCancelCmd)

	// Job submit flags
	jobSubmitCmd.Flags().StringVarP(&jobName, "job-name", "J", "", "Job name")
	jobSubmitCmd.Flags().StringVarP(&jobAccount, "account", "A", "", "Account to charge resources to")
	jobSubmitCmd.Flags().StringVarP(&jobPartition, "partition", "p", "", "Partition to submit job to")
	jobSubmitCmd.Flags().StringVar(&jobQos, "qos", "", "Quality of service")
	jobSubmitCmd.Flags().StringVar(&jobComment, "comment", "", "Comment")
	jobSubmitCmd.Flags().StringVarP(&jobNodes, "nodes", "N", "1", "Number of nodes required")
	jobSubmitCmd.Flags().IntVarP(&jobTasks, "ntasks", "n", 1, "Number of tasks")
	jobSubmitCmd.Flags().IntVarP(&jobCpusPerTask, "cpus-per-task", "c", 1, "CPUs per task")
	jobSubmitCmd.Flags().StringVar(&jobMemPerCpu, "mem-per-cpu", "", "Memory per CPU")
	jobSubmitCmd.Flags().StringVar(&jobMemPerNode, "mem", "", "Memory per node")
	jobSubmitCmd.Flags().StringVarP(&jobTimeLimit, "time", "t", "", "Time limit")
	jobSubmitCmd.Flags().StringVar(&jobBeginTime, "begin", "", "Begin time")
	jobSubmitCmd.Flags().StringVar(&jobDeadline, "deadline", "", "Job deadline")
	jobSubmitCmd.Flags().StringVarP(&jobStdIn, "input", "i", "", "Standard input file")
	jobSubmitCmd.Flags().StringVarP(&jobStdOut, "output", "o", "", "Standard output file")
	jobSubmitCmd.Flags().StringVarP(&jobStdErr, "error", "e", "", "Standard error file")
	jobSubmitCmd.Flags().BoolVarP(&jobHold, "hold", "H", false, "Submit job in held state")
	jobSubmitCmd.Flags().BoolVar(&jobRequeue, "requeue", false, "Requeue job on failure")
	jobSubmitCmd.Flags().BoolVar(&jobKillOnNodeFail, "kill-on-node-fail", true, "Kill job if any node fails")
	jobSubmitCmd.Flags().StringVarP(&jobArray, "array", "a", "", "Job array indices")
	jobSubmitCmd.Flags().StringVarP(&jobDependency, "dependency", "d", "", "Job dependencies")
	jobSubmitCmd.Flags().StringVar(&jobMailType, "mail-type", "", "Mail events to notify (BEGIN, END, FAIL, REQUEUE, ALL)")
	jobSubmitCmd.Flags().StringVar(&jobMailUser, "mail-user", "", "User to receive mail notification")
	jobSubmitCmd.Flags().IntVar(&jobNice, "nice", 0, "Scheduling priority adjustment")
	jobSubmitCmd.Flags().StringVarP(&jobConstraint, "constraint", "C", "", "Node constraints")
	jobSubmitCmd.Flags().BoolVar(&jobX11, "x11", false, "Enable X11 forwarding")
	jobSubmitCmd.Flags().StringVar(&jobGres, "gres", "", "Generic consumable resources")
	jobSubmitCmd.Flags().StringVar(&jobTresPerJob, "tres-per-job", "", "TRES per job")
	jobSubmitCmd.Flags().StringVar(&jobTresPerNode, "tres-per-node", "", "TRES per node")
	jobSubmitCmd.Flags().StringVar(&jobTresPerSocket, "tres-per-socket", "", "TRES per socket")
	jobSubmitCmd.Flags().StringVar(&jobTresPerTask, "tres-per-task", "", "TRES per task")
	jobSubmitCmd.Flags().StringVar(&jobCpusPerTres, "cpus-per-gpu", "", "CPUs per GPU")
	jobSubmitCmd.Flags().StringVar(&jobMemPerTres, "mem-per-gpu", "", "Memory per GPU")
	jobSubmitCmd.Flags().StringVarP(&jobLicenses, "licenses", "L", "", "License specifications")
	jobSubmitCmd.Flags().StringSliceVarP(&jobClusters, "clusters", "M", nil, "Clusters to submit job to")
	jobSubmitCmd.Flags().StringVar(&jobReservation, "reservation", "", "Resource reservation name")
	jobSubmitCmd.Flags().IntVar(&jobPriority, "priority", 0, "Job priority")
	jobSubmitCmd.Flags().StringVarP(&jobChdir, "chdir", "D", "", "Working directory for the job")
	jobSubmitCmd.Flags().StringVar(&jobWorkdir, "workdir", "", "Working directory for the job (alias for --chdir)")
	jobSubmitCmd.Flags().StringVar(&jobWckey, "wckey", "", "Job's workload characterization key")
	jobSubmitCmd.Flags().BoolVar(&jobExclusive, "exclusive", false, "Allocate nodes exclusively")
	jobSubmitCmd.Flags().BoolVar(&jobShared, "share", false, "Share nodes with other jobs")
	jobSubmitCmd.Flags().BoolVar(&jobOversubscribe, "oversubscribe", false, "Allow job to oversubscribe resources")
	jobSubmitCmd.Flags().BoolVar(&jobContiguous, "contiguous", false, "Require contiguous nodes")
	jobSubmitCmd.Flags().IntVar(&jobCoreSpec, "core-spec", 0, "Count of specialized cores per node")
	jobSubmitCmd.Flags().IntVar(&jobThreadSpec, "thread-spec", 0, "Count of specialized threads per node")
	jobSubmitCmd.Flags().IntVar(&jobMinCpus, "mincpus", 0, "Minimum CPUs per node")
	jobSubmitCmd.Flags().IntVar(&jobMinNodes, "min-nodes", 0, "Minimum number of nodes")
	jobSubmitCmd.Flags().IntVar(&jobMaxNodes, "max-nodes", 0, "Maximum number of nodes")
	jobSubmitCmd.Flags().IntVar(&jobSocketsPerNode, "sockets-per-node", 0, "Sockets per node")
	jobSubmitCmd.Flags().IntVar(&jobCoresPerSocket, "cores-per-socket", 0, "Cores per socket")
	jobSubmitCmd.Flags().IntVar(&jobThreadsPerCore, "threads-per-core", 0, "Threads per core")
	jobSubmitCmd.Flags().IntVar(&jobNtasksPerNode, "ntasks-per-node", 0, "Number of tasks per node")
	jobSubmitCmd.Flags().IntVar(&jobNtasksPerSocket, "ntasks-per-socket", 0, "Number of tasks per socket")
	jobSubmitCmd.Flags().IntVar(&jobNtasksPerCore, "ntasks-per-core", 0, "Number of tasks per core")
	jobSubmitCmd.Flags().StringVar(&jobNtasksPerTres, "ntasks-per-gpu", "", "Number of tasks per GPU")
	jobSubmitCmd.Flags().StringSliceVar(&jobEnvironment, "export", nil, "Environment variables to export")
	jobSubmitCmd.Flags().StringVar(&jobBurstBuffer, "bb", "", "Burst buffer specification")
	jobSubmitCmd.Flags().IntVar(&jobDelayBoot, "delay-boot", 0, "Delay boot for this number of minutes")
	jobSubmitCmd.Flags().StringVar(&jobNetwork, "network", "", "Network specification")

	// Job list flags
	jobListCmd.Flags().StringVarP(&jobUser, "user", "u", "", "Filter jobs by user")
	jobListCmd.Flags().StringVarP(&jobPartitionFilter, "partition", "p", "", "Filter jobs by partition")
	jobListCmd.Flags().StringVar(&jobStateFilter, "state", "", "Filter jobs by state")
	jobListCmd.Flags().StringVarP(&jobFormat, "format", "o", "", "Output format specifier")
}
