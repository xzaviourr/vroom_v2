package main

type ReqQueue struct {
	readyQueue      map[string][]*FuncReq // Task Id -> Queue mapping
	blockedQueue    map[string][]*FuncReq // Task Id -> Queue mapping
	resourceManager *ResourceManager      // Link to the resource manager
}

func (q *ReqQueue) Enque(funcReq *FuncReq) {
	runningInstances := q.resourceManager.taskStore.Instances[funcReq.TaskIdentifier]

	for _, instance := range runningInstances {
		if instance.Variant.Accuracy >= funcReq.Accuracy {
			if queue, ok := q.readyQueue[funcReq.TaskIdentifier]; ok {
				q.readyQueue[funcReq.TaskIdentifier] = append(queue, funcReq)
			} else {
				q.readyQueue[funcReq.TaskIdentifier] = []*FuncReq{funcReq}
			}
			funcReq.State = "ready"
			return
		}
	}

	// No valid instance is running
	if queue, ok := q.blockedQueue[funcReq.TaskIdentifier]; ok {
		q.blockedQueue[funcReq.TaskIdentifier] = append(queue, funcReq)
	} else {
		q.blockedQueue[funcReq.TaskIdentifier] = []*FuncReq{funcReq}
	}
	funcReq.State = "blocked"
}

func (q *ReqQueue) Front(taskIdentifier string, queueType int) *FuncReq {
	if queueType == 0 {
		return q.readyQueue[taskIdentifier][0]
	} else {
		return q.blockedQueue[taskIdentifier][0]
	}
}

func (q *ReqQueue) Deque(taskIdentifier string, queueType int) {
	if queueType == 0 {
		if len(q.readyQueue[taskIdentifier]) != 0 {
			q.readyQueue[taskIdentifier] = q.readyQueue[taskIdentifier][1:]
		}
	} else {
		if len(q.blockedQueue[taskIdentifier]) != 0 {
			q.blockedQueue[taskIdentifier] = q.blockedQueue[taskIdentifier][1:]
		}
	}
}

// func (q *ReqQueue) schedulingPolicy(k8s *kubernetes.Clientset, gpuResources map[string]GpuResource, resourceMutex *sync.Mutex) {
// 	for {
// 		time.Sleep(1 * time.Second)
// 		if len(q.items) > 0 {
// 			funcReq, _ := q.deque()

// 			current_ts := time.Now()
// 			// Time remaining before SLO miss
// 			remaining_time := funcReq.Deadline - float32(current_ts.Sub(funcReq.Timestamp)/time.Millisecond)

// 			variants, _ := getVariantsForReq(funcReq, remaining_time)
// 			var selected_variant FuncInfo

// 			if len(variants) == 0 {
// 				fmt.Println("SLO Miss")
// 				selected_variant, _ = getMinimumLatencyVariantForReq(funcReq)
// 			} else {
// 				selected_variant = variants[0]
// 			}
// 			// Variant selection logic
// 			// Node selection logic

// 			selected_node := "ub-10"

// 			if gpuResources[selected_node].VcoreAllocatable < selected_variant.GpuCores ||
// 				gpuResources[selected_node].VmemAllocatable < selected_variant.GpuMemory {
// 				q.addToQueue(funcReq)
// 			} else {
// 				// Update the availability of the resource
// 				resourceMutex.Lock()
// 				node_resource := gpuResources[selected_node]
// 				node_resource.VcoreAllocatable -= selected_variant.GpuCores
// 				node_resource.VmemAllocatable -= selected_variant.GpuMemory
// 				gpuResources[selected_node] = node_resource
// 				resourceMutex.Unlock()

// 				// Deploy the function
// 				deployFunc(selected_variant, k8s, funcReq.Uid)
// 			}
// 		}
// 		fmt.Println("Vcore : ", gpuResources["ub-10"].VcoreAllocatable)
// 		fmt.Println("Vmem : ", gpuResources["ub-10"].VmemAllocatable)
// 	}
// }
