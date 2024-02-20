package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
)

func (q *Queue) addToQueue(funcInfo FuncInfo) {
	fmt.Printf("Pod added to queue: %s", (funcInfo.task_identifier + funcInfo.variant_id))
	q.items = append(q.items, funcInfo)
}

func (q *Queue) deque() (FuncInfo, error) {
	if len(q.items) == 0 {
		return FuncInfo{}, fmt.Errorf("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *Queue) schedulingPolicy(k8s *kubernetes.Clientset) {
	for {
		time.Sleep(1 * time.Second)
		if len(q.items) > 0 {
			funcInfo, _ := q.deque()
			deployFunc(funcInfo, k8s)
		}
	}
}
