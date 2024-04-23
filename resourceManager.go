package main

import "fmt"

type ResourceManager struct {
	VariantStore  *VariantStore
	NodeStore     *NodeStore
	InstanceStore *InstanceStore
	TaskStore     *TaskStore
	RequestStore  *RequestStore
}

func initResourceManager() *ResourceManager {
	resourceManager := ResourceManager{
		VariantStore:  initVariantStore(),
		NodeStore:     initNodeStore(),
		InstanceStore: initInstanceStore(),
		TaskStore:     initTaskStore(),
		RequestStore:  initRequestStore(),
	}
	fmt.Println("Resource Manager Initialized Successfully")
	return &resourceManager
}
