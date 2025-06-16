package models

// JobSubmitRequest represents a job submission request.
type JobSubmitRequest struct {
	Script string    `json:"script,omitempty"`
	Job    *JobDescr `json:"job"`
}

// JobSubmitResponse represents a job submission response.
type JobSubmitResponse struct {
	Errors []string `json:"errors,omitempty"`
	JobID  int      `json:"job_id,omitempty"`
	StepID int      `json:"step_id,omitempty"`
}

// JobDescr represents a job description.
type JobDescr struct {
	Name                  string   `json:"name,omitempty"`
	Account               string   `json:"account,omitempty"`
	Partition             string   `json:"partition,omitempty"`
	QOS                   string   `json:"qos,omitempty"`
	Comment               string   `json:"comment,omitempty"`
	Nodes                 string    `json:"nodes,omitempty"`
	Tasks                 int      `json:"tasks,omitempty"`
	CPUsPerTask           int      `json:"cpus_per_task,omitempty"`
	MemPerCPU             string   `json:"mem_per_cpu,omitempty"`
	MemPerNode            string   `json:"mem_per_node,omitempty"`
	TimeLimit             string   `json:"time_limit,omitempty"`
	BeginTime             string   `json:"begin_time,omitempty"`
	Deadline              string   `json:"deadline,omitempty"`
	StdIn                 string   `json:"std_in,omitempty"`
	StdOut                string   `json:"std_out,omitempty"`
	StdErr                string   `json:"std_err,omitempty"`
	Hold                  bool     `json:"hold,omitempty"`
	Requeue               bool     `json:"requeue,omitempty"`
	KillOnNodeFail        bool     `json:"kill_on_node_fail,omitempty"`
	ArrayInx              string   `json:"array_inx,omitempty"`
	Dependency            string   `json:"dependency,omitempty"`
	MailType              string   `json:"mail_type,omitempty"`
	MailUser              string   `json:"mail_user,omitempty"`
	Nice                  int      `json:"nice,omitempty"`
	Constraints           string   `json:"constraints,omitempty"`
	X11                   bool     `json:"x11,omitempty"`
	GRES                  string   `json:"gres,omitempty"`
	TRESPerJob            string   `json:"tres_per_job,omitempty"`
	TRESPerNode           string   `json:"tres_per_node,omitempty"`
	TRESPerSocket         string   `json:"tres_per_socket,omitempty"`
	TRESPerTask           string   `json:"tres_per_task,omitempty"`
	CPUsPerTRES           string   `json:"cpus_per_tres,omitempty"`
	MemPerTRES            string   `json:"mem_per_tres,omitempty"`
	Licenses              string   `json:"licenses,omitempty"`
	Clusters              []string `json:"clusters,omitempty"`
	Reservation           string   `json:"reservation,omitempty"`
	Priority              int      `json:"priority,omitempty"`
	CurrentWorkingDir     string   `json:"current_working_directory,omitempty"`
	Workdir               string   `json:"workdir,omitempty"`
	WCKey                 string   `json:"wckey,omitempty"`
	Exclusive             bool     `json:"exclusive,omitempty"`
	Shared                bool     `json:"shared,omitempty"`
	Oversubscribe         bool     `json:"oversubscribe,omitempty"`
	Contiguous            bool     `json:"contiguous,omitempty"`
	CoreSpec              int      `json:"core_spec,omitempty"`
	ThreadSpec            int      `json:"thread_spec,omitempty"`
	MinCPUs               int      `json:"min_cpus,omitempty"`
	MinNodes              int      `json:"min_nodes,omitempty"`
	MaxNodes              int      `json:"max_nodes,omitempty"`
	SocketsPerNode        int      `json:"sockets_per_node,omitempty"`
	CoresPerSocket        int      `json:"cores_per_socket,omitempty"`
	ThreadsPerCore        int      `json:"threads_per_core,omitempty"`
	NTasksPerNode         int      `json:"ntasks_per_node,omitempty"`
	NTasksPerSocket       int      `json:"ntasks_per_socket,omitempty"`
	NTasksPerCore         int      `json:"ntasks_per_core,omitempty"`
	NTasksPerTRES         string   `json:"ntasks_per_tres,omitempty"`
	Environment           []string `json:"environment,omitempty"`
	BurstBuffer           string   `json:"burst_buffer,omitempty"`
	DelayBoot             int      `json:"delay_boot,omitempty"`
	Network               string   `json:"network,omitempty"`
	Script                string   `json:"script,omitempty"`
}

// JobInfo represents information about a job.
type JobInfo struct {
	Account               string   `json:"account"`
	AccrueTime            int64    `json:"accrue_time"`
	AdminComment          string   `json:"admin_comment"`
	AllocNode             string   `json:"alloc_node"`
	AllocSID              int      `json:"alloc_sid"`
	ArrayJobID            int      `json:"array_job_id"`
	ArrayTaskID           int      `json:"array_task_id"`
	ArrayMaxTasks         int      `json:"array_max_tasks"`
	ArrayTaskString       string   `json:"array_task_string"`
	AssocID               int      `json:"assoc_id"`
	BatchFeatures         string   `json:"batch_features"`
	BatchFlag             int      `json:"batch_flag"`
	BatchHost             string   `json:"batch_host"`
	Flags                 []string `json:"flags"`
	BurstBuffer           string   `json:"burst_buffer"`
	BurstBufferState      string   `json:"burst_buffer_state"`
	Cluster               string   `json:"cluster"`
	ClusterFeatures       string   `json:"cluster_features"`
	Command               string   `json:"command"`
	Comment               string   `json:"comment"`
	Contiguous            int      `json:"contiguous"`
	CoreSpec              int      `json:"core_spec"`
	ThreadSpec            int      `json:"thread_spec"`
	CPUsPerTask           int      `json:"cpus_per_task"`
	Dependency            string   `json:"dependency"`
	Derived               int      `json:"derived_ec"`
	EligibleTime          int64    `json:"eligible_time"`
	EndTime               int64    `json:"end_time"`
	ExcludedNodes         string   `json:"excluded_nodes"`
	ExitCode              int      `json:"exit_code"`
	Features              string   `json:"features"`
	FedOriginStr          string   `json:"fed_origin_str"`
	FedSiblingStr         string   `json:"fed_siblings_active_str"`
	FedSiblingViable      string   `json:"fed_siblings_viable_str"`
	GroupID               int      `json:"group_id"`
	JobID                 int      `json:"job_id"`
	JobResources          string   `json:"job_resources"`
	JobState              string   `json:"job_state"`
	LastSchedEval         int64    `json:"last_sched_evaluation"`
	Licenses              string   `json:"licenses"`
	MaxCPUs               int      `json:"max_cpus"`
	MaxNodes              int      `json:"max_nodes"`
	MCSLabel              string   `json:"mcs_label"`
	MemoryPerNode         int64    `json:"memory_per_node"`
	MemoryPerTask         int64    `json:"memory_per_task"`
	Name                  string   `json:"name"`
	Network               string   `json:"network"`
	Nice                  int      `json:"nice"`
	NodeCount             int      `json:"node_count"`
	Nodes                 string   `json:"nodes"`
	Partition             string   `json:"partition"`
	PreemptTime           int64    `json:"preempt_time"`
	PreSusTime            int64    `json:"pre_sus_time"`
	Priority              int      `json:"priority"`
	Profile               string   `json:"profile"`
	QOS                   string   `json:"qos"`
	Reboot                int      `json:"reboot"`
	ReqNodeList           string   `json:"req_nodes"`
	ReqSwitch             int      `json:"req_switch"`
	Requeue               int      `json:"requeue"`
	Reservation           string   `json:"reservation"`
	ResizeTime            int64    `json:"resize_time"`
	RestartCnt            int      `json:"restart_cnt"`
	SchedNodes            string   `json:"sched_nodes"`
	Shared                int      `json:"shared"`
	ShowFlags             []string `json:"show_flags"`
	SiteFactor            int      `json:"site_factor"`
	SocketsPerNode        int      `json:"sockets_per_node"`
	StartTime             int64    `json:"start_time"`
	StateDesc             string   `json:"state_description"`
	StateReason           string   `json:"state_reason"`
	StdErr                string   `json:"standard_error"`
	StdIn                 string   `json:"standard_input"`
	StdOut                string   `json:"standard_output"`
	SubmitTime            int64    `json:"submit_time"`
	SuspendTime           int64    `json:"suspend_time"`
	SystemComment         string   `json:"system_comment"`
	TimeLimit             int      `json:"time_limit"`
	TimeMin               int      `json:"time_min"`
	ThreadsPerCore        int      `json:"threads_per_core"`
	TRES                  string   `json:"tres_req_str"`
	TRESAlloc             string   `json:"tres_alloc_str"`
	UserID                int      `json:"user_id"`
	UserName              string   `json:"user_name"`
	WCKey                 string   `json:"wckey"`
	WorkDir               string   `json:"work_dir"`
}

// JobsResponse represents a response containing a list of jobs.
type JobsResponse struct {
	Errors []string  `json:"errors,omitempty"`
	Jobs   []JobInfo `json:"jobs,omitempty"`
}

// JobResponse represents a response containing a single job.
type JobResponse struct {
	Errors []string `json:"errors,omitempty"`
	Job    JobInfo  `json:"job,omitempty"`
}
