package main

import (
	"fmt"
	"sync"
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

func (q *Queue) schedulingPolicy(k8s *kubernetes.Clientset, gpuResources map[string]GpuResource, resourceMutex *sync.Mutex) {
	for {
		time.Sleep(1 * time.Second)
		if len(q.items) > 0 {
			funcReq, _ := q.deque()

			current_ts := time.Now()
			// Time remaining before SLO miss
			remaining_time := funcReq.deadline - float32(current_ts.Sub(funcReq.timestamp)/time.Millisecond)

			variants, _ := getVariantsForReq(funcReq, remaining_time)
			var selected_variant FuncInfo

			if len(variants) == 0 {
				fmt.Println("SLO Miss")
				selected_variant, _ = getMinimumLatencyVariantForReq(funcReq)
			} else {
				selected_variant = variants[0]
			}
			// Variant selection logic
			// Node selection logic

			selected_node := "ub-10"

			if gpuResources[selected_node].vcore_allocatable < selected_variant.gpu_cores || gpuResources[selected_node].vmem_allocatable < selected_variant.gpu_memory {
				q.addToQueue(funcReq)
			} else {
				// Update the availability of the resource
				resourceMutex.Lock()
				node_resource := gpuResources[selected_node]
				node_resource.vcore_allocatable -= selected_variant.gpu_cores
				node_resource.vmem_allocatable -= selected_variant.gpu_memory
				gpuResources[selected_node] = node_resource
				resourceMutex.Unlock()

				// Deploy the function
				deployFunc(selected_variant, k8s, funcReq.uid)
			}
		}
		fmt.Println("Vcore : ", gpuResources["ub-10"].vcore_allocatable)
		fmt.Println("Vmem : ", gpuResources["ub-10"].vmem_allocatable)
	}
}
