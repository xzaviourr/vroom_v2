package main

import (
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	Name             string // Node name in Kubernetes
	IpAddress        string // IP address
	GpuType          string // GPU model type
	VmemCapacity     int64  // Number of Vmem devices
	VcoreCapacity    int64  // Number of Vcore devices
	VmemAllocatable  int64  // Available Vmem devices
	VcoreAllocatable int64  // Available Vcore devices

	GpuMemoryUsage float32 // GPU memory usage in last monitoring cycle
	GpuCoreUsage   float32 // GPU core usage in last monitoring cycle

	RunningInstances      []*Instance // Pointers to running instances on this node
	RunningInstancesMutex sync.Mutex

	RequestCounter      int64 // Counter to track requests in one monitor cycle
	RequestCounterMutex sync.Mutex
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

func (n *Node) addRunningInstance(instance *Instance) {
	n.RunningInstancesMutex.Lock()
	n.RunningInstances = append(n.RunningInstances, instance)
	n.RunningInstancesMutex.Unlock()
}

func (n *Node) removeRunningInstance(instance *Instance) {
	n.RunningInstancesMutex.Lock()
	instances := n.RunningInstances

	indexToRemove := -1
	for ind, inst := range instances {
		if inst.Id == instance.Id {
			indexToRemove = ind
			break
		}
	}
	n.RunningInstances = append(instances[:indexToRemove], instances[indexToRemove+1:]...)
	n.RunningInstancesMutex.Unlock()
}

func (n *Node) newRequest() {
	n.RequestCounterMutex.Lock()
	n.RequestCounter += 1
	n.RequestCounterMutex.Unlock()
}

func (n *Node) resetRequestCounter() int64 {
	n.RequestCounterMutex.Lock()
	counter := n.RequestCounter
	n.RequestCounter = 0
	n.RequestCounterMutex.Unlock()
	return counter
}

type NodeStore struct {
	Nodes map[string]*Node // Node name -> Node Info structure
}

func initNodeStore() *NodeStore {
	nodeStore := NodeStore{
		Nodes: make(map[string]*Node),
	}
	return &nodeStore
}

func (ns *NodeStore) String() string {
	result := "Node Store\n"
	for _, node := range ns.Nodes {
		result += node.String() + "\n\n"
	}
	return result
}
