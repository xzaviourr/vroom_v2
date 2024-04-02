package main

// type Instance struct {
// 	Id       string
// 	funcInfo *FuncInfo
// 	port     int
// }

// type Node struct {
// 	Instances []*FuncInfo
// }

// type Variant struct {
// 	Nodes map[string]*Node // NodeId -> Node
// }

// type Task struct {
// 	RunningVariants map[string]*Variant // variantId -> Variant
// }

// type TaskManager struct {
// 	Tasks map[string]*Task // taskIdentifier -> Task
// }

// func (tm *TaskManager) addNewVariant(taskIdentifier string, variantId string, nodeId string, funcInfo *FuncInfo) {
// 	// Ensure the map is initialized
// 	if tm.Tasks == nil {
// 		tm.Tasks = make(map[string]*Task)
// 	}

// 	// Ensure the task exists or initialize it
// 	task, taskExists := tm.Tasks[taskIdentifier]
// 	if !taskExists {
// 		task = &Task{RunningVariants: make(map[string]*Variant)}
// 		tm.Tasks[taskIdentifier] = task
// 	}

// 	// Ensure the variant exists or initialize it
// 	variant, variantExists := task.RunningVariants[variantId]
// 	if !variantExists {
// 		variant = &Variant{Nodes: make(map[string]*Node)}
// 		task.RunningVariants[variantId] = variant
// 	}

// 	// Ensure the node exists or initialize it
// 	node, nodeExists := variant.Nodes[nodeId]
// 	if nodeExists {
// 		node.Instances = append(node.Instances, funcInfo)
// 	} else {
// 		node := &Node{Instances: []*FuncInfo{funcInfo}}
// 		variant.Nodes[nodeId] = node
// 	}
// }

// func (tm *TaskManager) removeVariant(taskIdentifier string, variantId string, nodeId string, funcInfo *FuncInfo) {
// 	tm.Tasks[taskIdentifier].RunningVariants[variantId].Nodes[nodeId].NumberOfInstances--
// }

// func getLiveTaskInfo(taskId string) {

// }
