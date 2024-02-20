package main

import (
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
)

type Queue struct {
	items []FuncReq
}

func (q *Queue) addToQueue(funcReq FuncReq) {
	fmt.Printf("Request %s added to pending queue", funcReq.uid)
	q.items = append(q.items, funcReq)
}

func (q *Queue) deque() (FuncReq, error) {
	if len(q.items) == 0 {
		return FuncReq{}, fmt.Errorf("pending queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *Queue) schedulingPolicy(k8s *kubernetes.Clientset) {
	for {
		time.Sleep(1 * time.Second)
		if len(q.items) > 0 {
			funcReq, _ := q.deque()
			variants, _ := getVariantsForReq(funcReq)
			deployFunc(variants[0], k8s)
		}
	}
}
