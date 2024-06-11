package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
)

type ReqQueue struct {
	ReadyQueue      map[string][]*FuncReq // Task Id -> Queue mapping
	BlockedQueue    map[string][]*FuncReq // Task Id -> Queue mapping
	ResourceManager *ResourceManager      // Link to the resource manager
	LoadBalancer    *LoadBalancer         // Link to the load balancer
}

func initReqQueue(resourceManager *ResourceManager, loadBalancer *LoadBalancer) *ReqQueue {
	reqQueue := ReqQueue{
		ReadyQueue:      make(map[string][]*FuncReq),
		BlockedQueue:    make(map[string][]*FuncReq),
		ResourceManager: resourceManager,
		LoadBalancer:    loadBalancer,
	}
	return &reqQueue
}

func (q *ReqQueue) Enque(funcReq *FuncReq) {
	runningInstances := q.ResourceManager.TaskStore.getInstances(funcReq.TaskIdentifier)

	// Add request to request store
	q.ResourceManager.RequestStore.newRequest(funcReq)

	// If valid instance is running, add request to ready queue
	for _, instance := range runningInstances {
		if instance.Variant.Accuracy >= funcReq.Accuracy {
			if queue, ok := q.ReadyQueue[funcReq.TaskIdentifier]; ok {
				q.ReadyQueue[funcReq.TaskIdentifier] = append(queue, funcReq)
			} else {
				q.ReadyQueue[funcReq.TaskIdentifier] = []*FuncReq{funcReq}
			}
			funcReq.State = "ready"
			return
		}
	}

	// No valid instance is running, add request to blocked queue
	if queue, ok := q.BlockedQueue[funcReq.TaskIdentifier]; ok {
		q.BlockedQueue[funcReq.TaskIdentifier] = append(queue, funcReq)
	} else {
		q.BlockedQueue[funcReq.TaskIdentifier] = []*FuncReq{funcReq}
	}
	funcReq.State = "blocked"
}

func (q *ReqQueue) Front(taskIdentifier string, queueType int) *FuncReq {
	if queueType == 0 {
		return q.ReadyQueue[taskIdentifier][0]
	} else {
		return q.BlockedQueue[taskIdentifier][0]
	}
}

func (q *ReqQueue) Deque(taskIdentifier string, queueType int) {
	if queueType == 0 {
		if len(q.ReadyQueue[taskIdentifier]) != 0 {
			q.ReadyQueue[taskIdentifier] = q.ReadyQueue[taskIdentifier][1:]
		}
	} else {
		if len(q.BlockedQueue[taskIdentifier]) != 0 {
			q.BlockedQueue[taskIdentifier] = q.BlockedQueue[taskIdentifier][1:]
		}
	}
}

func (q *ReqQueue) blockedQueueScheduler(resourceManager *ResourceManager) {
	fmt.Println("Blocked queue scheduler is running")
	for {
		time.Sleep(1 * time.Second)
		for taskId, reqList := range q.BlockedQueue {
			runningInstances := q.ResourceManager.TaskStore.getInstances(taskId)

			if len(reqList) == 0 {
				delete(q.BlockedQueue, taskId)
				continue
			}

			funcReq := q.Front(taskId, 1)

			for _, instance := range runningInstances {
				if instance.Variant.Accuracy >= funcReq.Accuracy {
					if queue, ok := q.ReadyQueue[funcReq.TaskIdentifier]; ok {
						q.ReadyQueue[funcReq.TaskIdentifier] = append(queue, funcReq)
					} else {
						q.ReadyQueue[funcReq.TaskIdentifier] = []*FuncReq{funcReq}
					}
					funcReq.State = "ready"
					q.Deque(taskId, 1)
				}
			}
		}
	}
}

func (q *ReqQueue) schedulingPolicy(clientset *kubernetes.Clientset, resourceManager *ResourceManager) {
	fmt.Println("Ready queue scheduler is running")
	for {
		time.Sleep(5 * time.Second)

		for taskIdentifier, slice := range q.ReadyQueue {
			instances := resourceManager.TaskStore.getInstances(taskIdentifier)
			if len(instances) == 0 {
				continue
			}

			numberOfInstances := len(instances)
			requestFlag := 0
			for _, req := range slice {
				go dispatch(req, instances[requestFlag].Url, instances[requestFlag].Id, instances[requestFlag].Variant.Accuracy, resourceManager.Logger)
				q.Deque(taskIdentifier, 0)
				requestFlag += 1
				requestFlag %= numberOfInstances
			}
		}
	}
}
