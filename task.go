package main

import (
	"fmt"
	"sync"
)

type TaskStore struct {
	Instances      map[string][]*Instance
	MaxLoadLimit   map[string]float32
	InstancesMutex map[string]*sync.Mutex
}

func initTaskStore() *TaskStore {
	taskStore := TaskStore{
		Instances:      make(map[string][]*Instance),
		InstancesMutex: make(map[string]*sync.Mutex),
		MaxLoadLimit:   make(map[string]float32),
	}
	return &taskStore
}

func (ts *TaskStore) getInstances(taskId string) []*Instance {
	if _, ok := ts.InstancesMutex[taskId]; !ok {
		ts.InstancesMutex[taskId] = &sync.Mutex{}
	}

	ts.InstancesMutex[taskId].Lock()
	instances := ts.Instances[taskId]
	ts.InstancesMutex[taskId].Unlock()
	return instances
}

func (ts *TaskStore) addInstance(instance *Instance) {
	taskId := instance.Variant.TaskId
	if _, ok := ts.InstancesMutex[taskId]; !ok {
		ts.InstancesMutex[taskId] = &sync.Mutex{}
	}

	ts.InstancesMutex[taskId].Lock()
	ts.MaxLoadLimit[taskId] += instance.Variant.Capacity
	ts.Instances[taskId] = append(ts.Instances[taskId], instance)
	ts.InstancesMutex[taskId].Unlock()
}

func (ts *TaskStore) deleteInstance(instance *Instance) {
	taskId := instance.Variant.TaskId
	ts.InstancesMutex[taskId].Lock()
	ts.MaxLoadLimit[taskId] -= instance.Variant.Capacity
	instances := ts.Instances[taskId]
	indexToRemove := -1
	for ind, inst := range instances {
		if inst.Id == instance.Id {
			indexToRemove = ind
			break
		}
	}
	ts.Instances[taskId] = append(instances[:indexToRemove], instances[indexToRemove+1:]...)
	ts.InstancesMutex[taskId].Unlock()
}

func (ts *TaskStore) getMaxLoadLimits() map[string]float32 {
	return ts.MaxLoadLimit
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
