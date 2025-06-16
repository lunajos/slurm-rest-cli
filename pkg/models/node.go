package models

// NodeInfo represents information about a node.
type NodeInfo struct {
	Architecture      string            `json:"architecture"`
	BurstBufferState  string            `json:"burst_buffer_state"`
	CPUs              int               `json:"cpus"`
	CPULoad           float64           `json:"cpu_load"`
	FreeMemory        int64             `json:"free_memory"`
	Gres              string            `json:"gres"`
	GresUsed          string            `json:"gres_used"`
	MCSLabel          string            `json:"mcs_label"`
	Name              string            `json:"name"`
	NextState         string            `json:"next_state"`
	Address           string            `json:"address"`
	BootTime          interface{}       `json:"boot_time"`
	Cores             int               `json:"cores"`
	CoresPerSocket    int               `json:"cores_per_socket"`
	Features          interface{}       `json:"features"`
	ActiveFeatures    interface{}       `json:"active_features"`
	NodeHostname      string            `json:"node_hostname"`
	NodeAddr          string            `json:"node_addr"`
	OperatingSystem   string            `json:"operating_system"`
	Owner             string            `json:"owner"`
	Partitions        []string          `json:"partitions"`
	Port              int               `json:"port"`
	RealMemory        int64             `json:"real_memory"`
	Reason            string            `json:"reason"`
	ReasonUID         int               `json:"reason_uid"`
	ReasonTime        int64             `json:"reason_time"`
	SlurmdStartTime   interface{}       `json:"slurmd_start_time"`
	Sockets           int               `json:"sockets"`
	State             interface{}       `json:"state"`
	ThreadsPerCore    int               `json:"threads_per_core"`
	TmpDisk           int64             `json:"tmp_disk"`
	Weight            int               `json:"weight"`
	Tres              string            `json:"tres"`
	Version           string            `json:"version"`
	Extra             interface{}       `json:"extra"`
	Comment           string            `json:"comment"`
	Energy            interface{}       `json:"energy"`
	CPUSpecList       string            `json:"cpu_spec_list"`
	CPUAlloc          int               `json:"cpu_alloc"`
	CPUErr            int               `json:"cpu_err"`
	CPUTot            int               `json:"cpu_tot"`
	CPUSpecCnt        int               `json:"cpu_spec_cnt"`
	MemSpecLimit      int64             `json:"mem_spec_limit"`
	FreeMemSpecLimit  int64             `json:"free_mem_spec_limit"`
	LastBusy          interface{}       `json:"last_busy"`
	Boards            int               `json:"boards"`
	PowerManagement   interface{}       `json:"power_mgmt"`
	PowerCap          int               `json:"power_cap"`
	ExtSensorsJoules  int64             `json:"ext_sensors_joules"`
	ExtSensorsWatts   int64             `json:"ext_sensors_watts"`
	ExtSensorsTemp    int               `json:"ext_sensors_temp"`
	ExtSensorsDataTime int64             `json:"ext_sensors_data_time"`
}

// NodesResponse represents a response containing a list of nodes.
type NodesResponse struct {
	Errors []string   `json:"errors,omitempty"`
	Nodes  []NodeInfo `json:"nodes,omitempty"`
}

// NodeResponse represents a response containing a single node.
type NodeResponse struct {
	Errors []string `json:"errors,omitempty"`
	Node   NodeInfo `json:"node,omitempty"`
}

// NodeUpdateRequest represents a node update request.
type NodeUpdateRequest struct {
	State    string `json:"state,omitempty"`
	Reason   string `json:"reason,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Features string `json:"features,omitempty"`
	Comment  string `json:"comment,omitempty"`
}
