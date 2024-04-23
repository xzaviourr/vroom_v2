package main

import "time"

func main() {
	resourceManager := initResourceManager()
	k8s := initKubernetes(resourceManager)
	loadBalancer := initLoadBalancer(k8s, resourceManager)
	reqQueue := initReqQueue(resourceManager, loadBalancer)

	go k8s.monitorPods()                                         // Serice that handles cleaning of pods
	go loadBalancer.monitorLoad()                                // Auto scalar of pods
	go reqQueue.blockedQueueScheduler(resourceManager)           // Blocked queue scheduler
	go reqQueue.schedulingPolicy(k8s.clientset, resourceManager) // Scheduler

	time.Sleep(2 * time.Second)
	router := initServer(reqQueue, resourceManager)
	_ = router.Run("0.0.0.0:8083")
}
