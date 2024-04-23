package main

import (
	"fmt"
	"sync"
	"time"
)

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

type RequestStore struct {
	Requests         map[string]*FuncReq // RequestId -> Request Info struct
	MeanIncomingLoad map[string]float32  // Mean Incoming load

	RequestCounter      map[string]int64 // Request counter for each task
	RequestCounterMutex map[string]*sync.Mutex
}

func initRequestStore() *RequestStore {
	requestStore := RequestStore{
		Requests:            make(map[string]*FuncReq),
		MeanIncomingLoad:    make(map[string]float32),
		RequestCounter:      make(map[string]int64),
		RequestCounterMutex: make(map[string]*sync.Mutex),
	}
	return &requestStore
}

func (rs *RequestStore) String() string {
	result := "Request Store\n"
	for _, request := range rs.Requests {
		result += request.String() + "\n\n"
	}
	return result
}

func (rs *RequestStore) newRequest(request *FuncReq) {
	taskId := request.TaskIdentifier
	if _, ok := rs.RequestCounterMutex[taskId]; !ok {
		rs.RequestCounterMutex[taskId] = &sync.Mutex{}
	}

	rs.RequestCounterMutex[taskId].Lock()
	rs.RequestCounter[taskId] += 1
	rs.Requests[taskId] = request
	rs.RequestCounterMutex[taskId].Unlock()
}

func (rs *RequestStore) getRequestCounter(taskId string) int64 {
	rs.RequestCounterMutex[taskId].Lock()
	counter := rs.RequestCounter[taskId]
	rs.RequestCounterMutex[taskId].Unlock()
	return counter
}

func (rs *RequestStore) resetRequestCounter(taskId string) int64 {
	rs.RequestCounterMutex[taskId].Lock()
	counter := rs.RequestCounter[taskId]
	rs.RequestCounter[taskId] = 0
	rs.RequestCounterMutex[taskId].Unlock()
	return counter
}
