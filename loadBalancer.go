package main

import (
	"fmt"
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
		time.Sleep(3 * time.Second) // Cycle duration

		taskCapacity := lb.ResourceManager.TaskStore.getMaxLoadLimits()
		for taskId := range lb.ResourceManager.RequestStore.Requests {
			activeLoad := lb.ResourceManager.RequestStore.resetRequestCounter(taskId)

			capacity, ok := taskCapacity[taskId]
			if !ok {
				capacity = 0
			}

			fmt.Println("Load Balancer : ", taskId, activeLoad, capacity)

			if float32(activeLoad) > capacity*0.8 { // If incoming load crossed 90% limit
				lb.createNewInstance(taskId, activeLoad, capacity)
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

func (lb *LoadBalancer) createNewInstance(taskId string, load int64, capacity float32) {
	// Variant Selection Logic
	variantId := "d3723a40-d95b-4d32-9d08-532770fb6ec2"
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
