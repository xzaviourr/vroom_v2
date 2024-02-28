package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func getNodeGpuReources(nodeName string, clientset *kubernetes.Clientset) GpuResource {
	// Get the node object
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	vmemCapacity := node.Status.Capacity["nvidia.com/vmem"]
	vcoreCapacity := node.Status.Capacity["nvidia.com/vcore"]
	vmemAllocatable := node.Status.Allocatable["nvidia.com/vmem"]
	vcoreAllocatable := node.Status.Allocatable["nvidia.com/vcore"]

	vmemCapacityInt, _ := vmemCapacity.AsInt64()
	vcoreCapacityInt, _ := vcoreCapacity.AsInt64()
	vmemAllocatableInt, _ := vmemAllocatable.AsInt64()
	vcoreAllocatableInt, _ := vcoreAllocatable.AsInt64()

	gpuResource := GpuResource{nodeName, int(vcoreCapacityInt), int(vmemCapacityInt), int(vcoreAllocatableInt), int(vmemAllocatableInt)}
	return gpuResource
}

func getGpuNodes(clientset *kubernetes.Clientset) []core.Node {
	// Specify the label selector
	labelSelector := "nos.nebuly.com/gpu-partitioning=mps"

	// Get the nodes with the specified label
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		panic(err.Error())
	}

	// Print the node names
	fmt.Println("Found these nodes compatible with the scheduler :")
	for _, node := range nodes.Items {
		fmt.Println("Node : ", node.Name)
	}
	return nodes.Items
}

func monitorPods(clientset *kubernetes.Clientset, gpuResources map[string]GpuResource, resourceMutex *sync.Mutex) {
	for {
		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
			LabelSelector: "service=vroom",
		})
		if err != nil {
			fmt.Println("Error getting pods:", err)
			continue
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase == core.PodSucceeded || pod.Status.Phase == core.PodFailed {
				alloted_resources := make(map[string]int)
				for key, value := range pod.Spec.Containers[0].Resources.Requests {
					alloted_resources[string(key)], _ = strconv.Atoi(value.String())
				}

				resourceMutex.Lock()
				node_resource := gpuResources[pod.Spec.NodeName]
				node_resource.vcore_allocatable += alloted_resources["nvidia.com/vcore"]
				node_resource.vmem_allocatable += alloted_resources["nvidia.com/vmem"]
				gpuResources[pod.Spec.NodeName] = node_resource
				resourceMutex.Unlock()

				// Pod has completed or failed, delete it
				err := clientset.CoreV1().Pods("default").Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				if err != nil {
					fmt.Printf("Error deleting pod %s: %v\n", pod.Name, err)
				} else {
					fmt.Printf("Pod %s deleted\n", pod.Name)
				}
			}
		}
		time.Sleep(5 * time.Second) // Adjust the polling interval as needed
	}
}

func createPodObject(funcInfo FuncInfo, task_id string) *core.Pod {
	fmt.Printf("Creating a new pod")
	podName := funcInfo.task_identifier + "-" + funcInfo.variant_id + "-" + task_id
	namespace := "default"
	imageName := funcInfo.image
	gpuMemory := resource.NewQuantity(int64(funcInfo.gpu_memory), resource.DecimalSI).DeepCopy()
	gpuCores := resource.NewQuantity(int64(funcInfo.gpu_cores), resource.DecimalSI).DeepCopy()
	user := int64(1000)
	batch_size := strconv.Itoa(funcInfo.batch_size)

	fmt.Printf("%s", podName)

	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"name":    funcInfo.task_identifier,
				"service": "vroom",
			},
		},
		Spec: core.PodSpec{
			HostIPC:       true,
			RestartPolicy: "OnFailure",
			SecurityContext: &core.PodSecurityContext{
				RunAsUser: &user,
			},
			Containers: []core.Container{
				{
					Name:            podName,
					Image:           imageName,
					ImagePullPolicy: core.PullIfNotPresent,
					Env: []core.EnvVar{
						{
							Name:  "batch_size",
							Value: batch_size,
						},
					},
					Resources: core.ResourceRequirements{
						Requests: map[core.ResourceName]resource.Quantity{
							"nvidia.com/vmem":  gpuMemory,
							"nvidia.com/vcore": gpuCores,
						},
						Limits: map[core.ResourceName]resource.Quantity{
							"nvidia.com/vmem":  gpuMemory,
							"nvidia.com/vcore": gpuCores,
						},
					},
				},
			},
		},
	}
}

func deployFunc(funcInfo FuncInfo, clientset *kubernetes.Clientset, task_id string) string {
	pod := createPodObject(funcInfo, task_id)
	podName := pod.Name
	_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	return podName
}
