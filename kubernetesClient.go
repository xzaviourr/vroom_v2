package main

import (
	"context"
	"fmt"
	"strconv"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	_ "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func initializeNodes(clientset *kubernetes.Clientset, resourceManager *ResourceManager) {
	labelSelector := "nos.nebuly.com/gpu-partitioning=mps"
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodes.Items {
		vmemCapacity := node.Status.Capacity["nvidia.com/vmem"]
		vcoreCapacity := node.Status.Capacity["nvidia.com/vcore"]
		vmemAllocatable := node.Status.Allocatable["nvidia.com/vmem"]
		vcoreAllocatable := node.Status.Allocatable["nvidia.com/vcore"]
		vmemCapacityInt, _ := vmemCapacity.AsInt64()
		vcoreCapacityInt, _ := vcoreCapacity.AsInt64()
		vmemAllocatableInt, _ := vmemAllocatable.AsInt64()
		vcoreAllocatableInt, _ := vcoreAllocatable.AsInt64()

		ipAddress := ""
		for _, address := range node.Status.Addresses {
			if address.Type == core.NodeInternalIP {
				ipAddress = address.Address
				break
			}
		}

		newNode := Node{
			Name:             node.Name,
			IpAddress:        ipAddress,
			GpuType:          node.Labels["gpu-type"],
			VmemCapacity:     vmemCapacityInt,
			VcoreCapacity:    vcoreCapacityInt,
			VmemAllocatable:  vmemAllocatableInt,
			VcoreAllocatable: vcoreAllocatableInt,
			RunningInstances: []*Instance{},
			GpuMemoryUsage:   0.0,
			GpuCoreUsage:     0.0,
		}
		resourceManager.nodeStore.Nodes[node.Name] = &newNode
	}
}

func createPodForInstance(instance *Instance) *core.Pod {
	podName := instance.Id
	namespace := "default"
	imageName := instance.Variant.Image
	interalPort := instance.Variant.Port
	gpuMemory := resource.NewQuantity(int64(instance.Variant.GpuMemory), resource.DecimalSI).DeepCopy()
	gpuCores := resource.NewQuantity(int64(instance.Variant.GpuCores), resource.DecimalSI).DeepCopy()
	user := int64(1000)
	batch_size := strconv.Itoa(int(instance.Variant.BatchSize))

	fmt.Printf("%s", podName)

	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"name":    instance.Variant.TaskId,
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
					Ports: []core.ContainerPort{
						{
							ContainerPort: int32(interalPort),
						},
					},
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

func createServiceForInstance(instance *Instance) *core.Service {
	serviceName := instance.Id + "-service" // Define a name for your service
	internalPort := instance.Variant.Port
	externalPort := instance.Port // Define the port your service will listen on
	namespace := "default"        // Define the namespace for your service

	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
		},
		Spec: core.ServiceSpec{
			Selector: map[string]string{
				"name": instance.Id, // Match labels with your pod
			},
			Ports: []core.ServicePort{
				{
					Protocol:   core.ProtocolTCP,
					Port:       int32(externalPort),
					TargetPort: intstr.FromInt(int(internalPort)),
				},
			},
		},
	}
}

func deployInstance(instance *Instance, clientset *kubernetes.Clientset) string {
	pod := createPodForInstance(instance)
	service := createServiceForInstance(instance)

	// Deploy the Pod
	_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error deploying Pod:", err)
		return ""
	}

	// Deploy the Service
	_, err = clientset.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error deploying Service:", err)
		return ""
	}

	return pod.Name
}

// func monitorPods(clientset *kubernetes.Clientset, gpuResources map[string]GpuResource, resourceMutex *sync.Mutex) {
// 	for {
// 		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
// 			LabelSelector: "service=vroom",
// 		})
// 		if err != nil {
// 			fmt.Println("Error getting pods:", err)
// 			continue
// 		}

// 		for _, pod := range pods.Items {
// 			if pod.Status.Phase == core.PodSucceeded || pod.Status.Phase == core.PodFailed {
// 				alloted_resources := make(map[string]int)
// 				for key, value := range pod.Spec.Containers[0].Resources.Requests {
// 					alloted_resources[string(key)], _ = strconv.Atoi(value.String())
// 				}

// 				resourceMutex.Lock()
// 				node_resource := gpuResources[pod.Spec.NodeName]
// 				node_resource.VcoreAllocatable += alloted_resources["nvidia.com/vcore"]
// 				node_resource.VmemAllocatable += alloted_resources["nvidia.com/vmem"]
// 				gpuResources[pod.Spec.NodeName] = node_resource
// 				resourceMutex.Unlock()

// 				// Pod has completed or failed, delete it
// 				err := clientset.CoreV1().Pods("default").Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
// 				if err != nil {
// 					fmt.Printf("Error deleting pod %s: %v\n", pod.Name, err)
// 				} else {
// 					fmt.Printf("Pod %s deleted\n", pod.Name)
// 				}
// 			}
// 		}
// 		time.Sleep(5 * time.Second) // Adjust the polling interval as needed
// 	}
// }
