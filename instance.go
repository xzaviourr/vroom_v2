package main

import (
	"fmt"
	"sync"
)

type Instance struct {
	Id      string   // Unique id for the instance
	Variant *Variant // Variant info of this instance
	Node    *Node    // Node info of this instance
	Port    int64    // Service Port exposed externally
	Url     string   // Service url : Node IP + External Port + Service End Point

	State      string // State (pending, running, peak, overload, failed)
	StateMutex sync.Mutex

	RequestCounter      int64      // Request arrived in the last monitor cycle
	RequestCounterMutex sync.Mutex // Request counter mutex
}

func (i *Instance) getState() string {
	i.StateMutex.Lock()
	state := i.State
	i.StateMutex.Unlock()
	return state
}

func (i *Instance) setState(state string) {
	i.StateMutex.Lock()
	i.State = state
	i.StateMutex.Unlock()
}

func (i *Instance) newRequest() {
	i.RequestCounterMutex.Lock()
	i.RequestCounter += 1
	i.RequestCounterMutex.Unlock()
}

func (i *Instance) resetRequestCounter() int64 {
	i.RequestCounterMutex.Lock()
	counter := i.RequestCounter
	i.RequestCounter = 0
	i.RequestCounterMutex.Unlock()
	return counter
}

func (i *Instance) String() string {
	return fmt.Sprintf(
		"Instance\nId: %s\nVariant Id: %s\nNode Name: %s\nPort: %d\nUrl: %s\nState: %s\n",
		i.Id, i.Variant.Id, i.Node.Name, i.Port, i.Url, i.State,
	)
}

type InstanceStore struct {
	Instances map[string]*Instance // InstanceId -> Instance info struct
}

func initInstanceStore() *InstanceStore {
	instanceStore := InstanceStore{
		Instances: make(map[string]*Instance),
	}
	return &instanceStore
}

func (is *InstanceStore) String() string {
	result := "Instance Store\n"
	for _, instance := range is.Instances {
		result += instance.String() + "\n\n"
	}
	return result
}
