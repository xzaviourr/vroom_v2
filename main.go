package main

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func initKubernetes() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"),
		"(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// build configuration from the config file.
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	clientset, _ := kubernetes.NewForConfig(config)
	return clientset
}

func main() {
	resourceManager := ResourceManager{
		variantStore: VariantStore{
			Variants: make(map[string]*Variant),
		},
		nodeStore: NodeStore{
			Nodes: make(map[string]*Node),
		},
		instanceStore: InstanceStore{
			Instances: make(map[string]*Instance),
		},
		taskStore: TaskStore{
			Instances: make(map[string][]*Instance),
		},
		requestStore: RequestStore{
			Requests: make(map[string]*FuncReq),
		},
	}

	reqQueue := ReqQueue{
		readyQueue:      make(map[string][]*FuncReq),
		blockedQueue:    make(map[string][]*FuncReq),
		resourceManager: &resourceManager,
	}

	k8s := initKubernetes()
	initializeNodes(k8s, &resourceManager)
	setupDb()
	initializeVariants(&resourceManager)

	// var resourceMutex sync.Mutex

	// go reqQueue.schedulingPolicy(k8s, gpuResources, &resourceMutex)
	// go monitorPods(k8s, gpuResources, &resourceMutex)

	router := initServer(&reqQueue, &resourceManager)
	_ = router.Run("0.0.0.0:8083")
}
