package main

import (
	"fmt"
	"strings"
	"time"
)

type Variant struct {
	Id             string  // Unique id for the variant
	TaskId         string  `json:"task-identifier"` // Unique task id
	GpuMemory      int64   `json:"gpu-memory"`      // GPU memory (in MB)
	GpuCores       int64   `json:"gpu-cores"`       // Percentage of GPU cores
	Image          string  `json:"image"`           // Docker image file
	StartupLatency float32 `json:"startup-latency"` // Time to start the container
	MinLatency     float32 `json:"min-latency"`     // Latency for request when load is minimum
	MeanLatency    float32 `json:"mean-latency"`    // Latency for request of when load is at average
	MaxLatency     float32 `json:"max-latency"`     // Latency for request when load is at its peak
	Accuracy       float32 `json:"accuracy"`        // Model accuracy
	BatchSize      int64   `json:"batch-size"`      // Batch size
	EndPoint       string  `json:"end-point"`       // Url endpoint extension to access the service
	Port           int64   `json:"port"`            // Port at which service is running internally
}

type VariantStore struct {
	Variants map[string]*Variant // VariantId -> Variant Info structure
}

type Node struct {
	Name             string      // Node name in Kubernetes
	IpAddress        string      // IP address
	GpuType          string      // GPU model type
	VmemCapacity     int64       // Number of Vmem devices
	VcoreCapacity    int64       // Number of Vcore devices
	VmemAllocatable  int64       // Available Vmem devices
	VcoreAllocatable int64       // Available Vcore devices
	RunningInstances []*Instance // Pointers to running instances on this node
	GpuMemoryUsage   float32     // GPU memory usage in last monitoring cycle
	GpuCoreUsage     float32     // GPU core usage in last monitoring cycle
}

type NodeStore struct {
	Nodes map[string]*Node // Node name -> Node Info structure
}

type Instance struct {
	Id      string   // Unique id for the instance
	Variant *Variant // Variant info of this instance
	Node    *Node    // Node info of this instance
	Port    int64    // Service Port exposed externally
	Url     string   // Service url : Node IP + External Port + Service End Point
	State   string   // State (initializing, ready, peak, overload, error)
}

type InstanceStore struct {
	Instances map[string]*Instance // InstanceId -> Instance info struct
}

type TaskStore struct {
	Instances map[string][]*Instance // TaskId -> List of Instance info struct
}

type FuncReq struct {
	Uid                string    // Unique id for this request
	TaskIdentifier     string    `json:"task-identifier"` // Requested task
	Deadline           float32   `json:"deadline"`        // Deadline (in milli second)
	Accuracy           float32   `json:"accuracy"`        // Accuracy requirement (in percentage)
	Args               string    `json:"args"`            // Arguments to be passed
	ResponseUrl        string    `json:"response-url"`    // Response to be sent to this url
	RequestSize        int       `json:"request-size"`    // Number of inference queries
	RegistrationTs     time.Time // Timestamp when this request got registered
	DeployInstanceTs   time.Time // Timestamp when instance deployment was initiated
	SentForExecutionTs time.Time // Timestamp when API call was made to execute this request
	ResponseTs         time.Time // Timestamp when the response was sent back
	SelectedNode       string    // Id of the node on which this request was executed
	State              string    // State (new, ready, blocked, running)
}

type RequestStore struct {
	Requests map[string]*FuncReq // RequestId -> Request Info struct
}

type ResourceManager struct {
	variantStore  VariantStore
	nodeStore     NodeStore
	instanceStore InstanceStore
	taskStore     TaskStore
	requestStore  RequestStore
}

// =============================================================================================================================
// String functions for all the entities

func (v *Variant) String() string {
	return fmt.Sprintf(
		"Function Variant\nId: %s\nTask Id: %s\nGPU Memory: %d\nGPU Cores: %d\nImage: %s\nStartup Latency: %f\n"+
			"Min Latency: %f\nMean Latency: %f\nMax Latency: %f\nAccuracy: %f\nBatch Size: %d\nEnd Point: %s\nPort: %d\n",
		v.Id, v.TaskId, v.GpuMemory, v.GpuCores, v.Image, v.StartupLatency,
		v.MinLatency, v.MeanLatency, v.MaxLatency, v.Accuracy, v.BatchSize, v.EndPoint, v.Port,
	)
}

func (vs *VariantStore) String() string {
	result := "Variant Store\n"
	for _, variant := range vs.Variants {
		result += variant.String() + "\n\n"
	}
	return result
}

func (i *Instance) String() string {
	return fmt.Sprintf(
		"Instance\nId: %s\nVariant Id: %s\nNode Name: %s\nPort: %d\nUrl: %s\nState: %s\n",
		i.Id, i.Variant.Id, i.Node.Name, i.Port, i.Url, i.State,
	)
}

func (is *InstanceStore) String() string {
	result := "Instance Store\n"
	for _, instance := range is.Instances {
		result += instance.String() + "\n\n"
	}
	return result
}

func (n *Node) String() string {
	var runningInstancesInfo []string
	for _, instance := range n.RunningInstances {
		runningInstancesInfo = append(runningInstancesInfo, instance.String())
	}

	return fmt.Sprintf(
		"Node\nName: %s\nIpAddress: %s\nGpuType: %s\nVmemCapacity: %d\nVcoreCapacity: %d\nVmemAllocatable: %d\n"+
			"VcoreAllocatable: %d\nRunning Instances:\n%sGpuMemoryUsage: %f\nGpuCoreUsage: %f\n",
		n.Name, n.IpAddress, n.GpuType, n.VmemCapacity, n.VcoreCapacity, n.VmemAllocatable, n.VcoreAllocatable,
		strings.Join(runningInstancesInfo, "\n"), n.GpuMemoryUsage, n.GpuCoreUsage,
	)
}

func (ns *NodeStore) String() string {
	result := "Node Store\n"
	for _, node := range ns.Nodes {
		result += node.String() + "\n\n"
	}
	return result
}

func (ts *TaskStore) String() string {
	result := "Task Store\n"
	for taskId, instances := range ts.Instances {
		result += fmt.Sprintf("Task ID: %s\n", taskId)
		for _, instance := range instances {
			result += instance.String() + "\n\n"
		}
	}
	return result
}

func (fr *FuncReq) String() string {
	return fmt.Sprintf(
		"FuncReq\nUid: %s\nTaskIdentifier: %s\nDeadline: %f\nAccuracy: %f\nArgs: %s\nResponseUrl: %s\n"+
			"RegistrationTs: %s\nDeployInstanceTs: %s\nSentForExecutionTs: %s\nResponseTs: %s\n"+
			"SelectedNode: %s\nReqestSize: %d\nState: %s\n",
		fr.Uid, fr.TaskIdentifier, fr.Deadline, fr.Accuracy, fr.Args, fr.ResponseUrl,
		fr.RegistrationTs.String(), fr.DeployInstanceTs.String(), fr.SentForExecutionTs.String(),
		fr.ResponseTs.String(), fr.SelectedNode, fr.RequestSize, fr.State,
	)
}

func (rs *RequestStore) String() string {
	result := "Request Store\n"
	for _, request := range rs.Requests {
		result += request.String() + "\n\n"
	}
	return result
}
