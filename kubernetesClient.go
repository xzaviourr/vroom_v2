package main

import (
	"context"
	"fmt"
	"strconv"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

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
				"name": funcInfo.task_identifier,
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
