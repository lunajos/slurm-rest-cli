package models

// PartitionInfo represents information about a partition.
type PartitionInfo struct {
	AllowGroups         string   `json:"allow_groups"`
	AllowAccounts       string   `json:"allow_accounts"`
	AllowQOS            string   `json:"allow_qos"`
	AllocNodes          string   `json:"alloc_nodes"`
	Alternate           string   `json:"alternate"`
	CPUBind             string   `json:"cpu_bind"`
	DefMemPerCPU        int64    `json:"def_mem_per_cpu"`
	DefMemPerNode       int64    `json:"def_mem_per_node"`
	DefaultTime         string   `json:"default_time"`
	DenyAccounts        string   `json:"deny_accounts"`
	DenyQOS             string   `json:"deny_qos"`
	DisableRootJobs     int      `json:"disable_root_jobs"`
	ExclusiveUser       int      `json:"exclusive_user"`
	Flags               []string `json:"flags"`
	GraceTime           int      `json:"grace_time"`
	Hidden              int      `json:"hidden"`
	MaxMemPerCPU        int64    `json:"max_mem_per_cpu"`
	MaxMemPerNode       int64    `json:"max_mem_per_node"`
	MaxNodes            int      `json:"max_nodes"`
	MaxTime             string   `json:"max_time"`
	MinNodes            int      `json:"min_nodes"`
	Name                string   `json:"name"`
	NodeList            string   `json:"nodes"`
	OverSubscribe       string   `json:"over_subscribe"`
	OverTimeLimit       int      `json:"over_time_limit"`
	PreemptMode         string   `json:"preempt_mode"`
	Priority            int      `json:"priority"`
	PriorityJobFactor   int      `json:"priority_job_factor"`
	PriorityTier        int      `json:"priority_tier"`
	QOS                 string   `json:"qos"`
	ReqResv             int      `json:"req_resv"`
	RootOnly            int      `json:"root_only"`
	SelectTypeParams    int      `json:"select_type_parameters"`
	State               string   `json:"state"`
	TotalCPUs           int      `json:"total_cpus"`
	TotalNodes          int      `json:"total_nodes"`
	Tres                string   `json:"tres"`
	MaxCPUsPerNode      int      `json:"max_cpus_per_node"`
	AllowAlloc          int      `json:"allow_alloc_nodes"`
	AllowAccounts2      []string `json:"allowed_accounts"`
	AllowGroups2        []string `json:"allowed_groups"`
	AllowQOS2           []string `json:"allowed_qos"`
	DenyAccounts2       []string `json:"denied_accounts"`
	DenyQOS2            []string `json:"denied_qos"`
	PreemptMode2        []string `json:"preempt_mode_string"`
}

// PartitionListResponse represents a response containing a list of partitions.
type PartitionListResponse struct {
	Errors     []string        `json:"errors,omitempty"`
	Partitions []interface{}   `json:"partitions,omitempty"`
	LastUpdate interface{}     `json:"last_update,omitempty"`
	Meta       interface{}     `json:"meta,omitempty"`
	Warnings   []interface{}   `json:"warnings,omitempty"`
}

// PartitionGetResponse represents a response containing a single partition.
type PartitionGetResponse struct {
	Errors    []string      `json:"errors,omitempty"`
	Partition PartitionInfo `json:"partition,omitempty"`
}

// PartitionUpdateRequest represents a request to update a partition.
type PartitionUpdateRequest struct {
	State               string   `json:"state,omitempty"`
	DefaultTime         string   `json:"default_time,omitempty"`
	MaxTime             string   `json:"max_time,omitempty"`
	Priority            int      `json:"priority,omitempty"`
	AllowGroups         string   `json:"allow_groups,omitempty"`
	AllowAccounts       string   `json:"allow_accounts,omitempty"`
	AllowQOS            string   `json:"allow_qos,omitempty"`
	DenyAccounts        string   `json:"deny_accounts,omitempty"`
	DenyQOS             string   `json:"deny_qos,omitempty"`
	OverSubscribe       string   `json:"over_subscribe,omitempty"`
	Hidden              int      `json:"hidden,omitempty"`
	MaxNodes            int      `json:"max_nodes,omitempty"`
	MinNodes            int      `json:"min_nodes,omitempty"`
	DefMemPerCPU        int64    `json:"def_mem_per_cpu,omitempty"`
	DefMemPerNode       int64    `json:"def_mem_per_node,omitempty"`
	MaxMemPerCPU        int64    `json:"max_mem_per_cpu,omitempty"`
	MaxMemPerNode       int64    `json:"max_mem_per_node,omitempty"`
	Nodes               string   `json:"nodes,omitempty"`
	AllocNodes          string   `json:"alloc_nodes,omitempty"`
	Alternate           string   `json:"alternate,omitempty"`
	GraceTime           int      `json:"grace_time,omitempty"`
	QOS                 string   `json:"qos,omitempty"`
	DisableRootJobs     int      `json:"disable_root_jobs,omitempty"`
	ExclusiveUser       int      `json:"exclusive_user,omitempty"`
	OverTimeLimit       int      `json:"over_time_limit,omitempty"`
	PreemptMode         string   `json:"preempt_mode,omitempty"`
}
