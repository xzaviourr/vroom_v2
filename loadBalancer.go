package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LoadBalancer struct {
	ResourceManager *ResourceManager
	k8s             *K8s
}

func initLoadBalancer(k8s *K8s, resourceManager *ResourceManager) *LoadBalancer {
	loadBalancer := LoadBalancer{
		ResourceManager: resourceManager,
		k8s:             k8s,
	}
	return &loadBalancer
}

func (lb *LoadBalancer) monitorLoad() {
	fmt.Println("Load Balancer is running")
	for {
		time.Sleep(5 * time.Second) // Cycle duration

		taskCapacity := lb.ResourceManager.TaskStore.getMaxLoadLimits()
		for taskId := range lb.ResourceManager.RequestStore.Requests {
			activeLoad := lb.ResourceManager.RequestStore.resetRequestCounter(taskId)

			capacity, ok := taskCapacity[taskId]
			if !ok {
				capacity = 0
			}

			fmt.Println("Load Balancer : ", taskId, float32(activeLoad)/5, float32(capacity))

			if (float32(activeLoad) / 5) > (float32(capacity) * 0.8) { // If incoming load crossed 90% limit
				lb.scaleOperation(taskId, float32(activeLoad)/5, capacity)
			}
		}
	}
}

func (lb *LoadBalancer) getFreePort(nodeName string) int64 {
	node := lb.ResourceManager.NodeStore.Nodes[nodeName]
	usedPorts := make(map[int64]bool)
	for _, instance := range node.RunningInstances {
		usedPorts[instance.Port] = true
	}

	lowPort := int64(49152)
	highPort := int64(65535)
	for port := lowPort; port < highPort; port++ {
		if !usedPorts[port] {
			return port
		}
	}
	return -1
}

func (lb *LoadBalancer) findPerformanceKneePointVariant(taskId string, currentArrivalRate float32) *Variant {
	resourceVariants := lb.ResourceManager.VariantStore.getTaskVariants(taskId)

	// Sort the variant list by (GPUmemory, GPUcores)
	sort.Slice(resourceVariants, func(i, j int) bool {
		if resourceVariants[i].GpuMemory == resourceVariants[j].GpuMemory {
			return resourceVariants[i].GpuCores < resourceVariants[j].GpuCores
		}
		return resourceVariants[i].GpuMemory < resourceVariants[j].GpuMemory
	})

	// Initialize knee_point and best_ratio
	var kneePoint *Variant
	bestRatio := 0.0

	// Iterate over each resource variant in the sorted list
	for _, variant := range resourceVariants {
		performanceRatio := float64(variant.Capacity) / float64((float32(variant.GpuMemory)/float32(16))+(float32(variant.GpuCores)/float32(100)))
		if performanceRatio > bestRatio {
			bestRatio = performanceRatio
			kneePoint = variant
		}
	}

	// Return the knee_point variant
	return kneePoint
}

func (lb *LoadBalancer) findResourceVariantGroup(taskId string, requiredCapacity float32, maxColocationFactor int) []*Variant {
	variants := lb.ResourceManager.VariantStore.getTaskVariants(taskId)
	selectedGroup := []*Variant{}
	minimumResources := float32(math.MaxFloat32)

	// Helper function to calculate total resources
	calculateTotalResources := func(group []*Variant) float32 {
		totalResources := float32(0)
		for _, v := range group {
			totalResources += float32(v.GpuMemory) * float32(v.GpuCores)
		}
		return totalResources
	}

	// Helper function to calculate total throughput
	calculateTotalThroughput := func(group []*Variant) float32 {
		totalThroughput := float32(0)
		for _, v := range group {
			totalThroughput += float32(v.Capacity)
		}
		return totalThroughput
	}

	// Generate all combinations of resource variants
	var generateCombinations func([]*Variant, int, int, []*Variant)
	generateCombinations = func(variants []*Variant, start, depth int, current []*Variant) {
		if depth == 0 {
			totalResources := calculateTotalResources(current)
			totalThroughput := calculateTotalThroughput(current)
			if totalThroughput >= requiredCapacity && totalResources < minimumResources {
				minimumResources = totalResources
				selectedGroup = make([]*Variant, len(current))
				copy(selectedGroup, current)
			}
			return
		}
		for i := start; i <= len(variants)-depth; i++ {
			generateCombinations(variants, i+1, depth-1, append(current, variants[i]))
		}
	}

	// Iterate over each colocation factor
	for cf := 1; cf <= maxColocationFactor; cf++ {
		generateCombinations(variants, 0, cf, []*Variant{})
	}

	return selectedGroup
}

func (lb *LoadBalancer) scaleOperation(taskId string, load float32, capacity float32) {
	// Variant Selection Logic
	if capacity == 0 { // First variant
		variant := lb.findPerformanceKneePointVariant(taskId, load)
		lb.createNewInstance(variant.Id)
	} else { // Scaling operation
		res := float32(load - capacity)
		variants := lb.findResourceVariantGroup(taskId, res, 4)
		for _, variant := range variants {
			lb.createNewInstance(variant.Id)
		}
	}
}

func (lb *LoadBalancer) createNewInstance(variantId string) {
	// variantId := "4788d252-3481-4618-a83f-87ed1bfb8875"
	// Node Selection Logic
	nodeName := "ub-10"

	instanceId := "v" + uuid.New().String()
	variant := lb.ResourceManager.VariantStore.Variants[variantId]
	node := lb.ResourceManager.NodeStore.Nodes[nodeName]
	port := lb.getFreePort(nodeName)
	serviceUrl := node.IpAddress + ":" + strconv.Itoa(int(port)) + variant.EndPoint

	instance := Instance{
		Id:                  instanceId,
		Variant:             variant,
		Node:                node,
		Port:                port,
		Url:                 serviceUrl,
		State:               "pending",
		StateMutex:          sync.Mutex{},
		RequestCounter:      0,
		RequestCounterMutex: sync.Mutex{},
	}
	lb.ResourceManager.InstanceStore.Instances[instanceId] = &instance
	lb.ResourceManager.TaskStore.addInstance(&instance)
	lb.ResourceManager.NodeStore.Nodes[nodeName].addRunningInstance(&instance)

	fmt.Println("New instance created : ", instance.String())
	lb.k8s.deployInstance(&instance)
}
