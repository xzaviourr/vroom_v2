package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type GpuResource struct {
	node_id           string
	vcore_capacity    int
	vmem_capacity     int
	vcore_allocatable int
	vmem_allocatable  int
}

type FuncInfo struct {
	variant_id      string
	task_identifier string
	gpu_memory      int
	gpu_cores       int
	image           string
	latency         float32
	accuracy        float32
	batch_size      int
}

type FuncReq struct {
	uid             string
	task_identifier string
	deadline        float32
	accuracy        float32
	timestamp       time.Time
}

func InitQueue() *Queue {
	queue := Queue{}
	return &queue
}

func InitResources(clientset *kubernetes.Clientset) map[string]GpuResource {
	gpu_nodes := getGpuNodes(clientset)
	resources := make(map[string]GpuResource)
	for _, node := range gpu_nodes {
		resources[node.Name] = getNodeGpuReources(node.Name, clientset)
	}
	return resources
}

func initKubernetes() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// build configuration from the config file.
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	clientset, _ := kubernetes.NewForConfig(config)
	return clientset
}

func initRouter(queue *Queue) *gin.Engine {
	r := gin.Default()
	r.GET("/run", func(c *gin.Context) {
		var funcReq FuncReq

		// Unique id for the new task
		funcReq.uid = uuid.New().String()

		// Task identifer
		funcReq.task_identifier = c.Query("task_id")

		// Deadline constraint for the task
		deadline64, err := strconv.ParseFloat(c.Query("deadline"), 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid deadline value"})
		}
		funcReq.deadline = float32(deadline64)

		// Accuracy constraint for the task
		accuracy64, err := strconv.ParseFloat(c.Query("accuracy"), 32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid accuracy value"})
		}
		funcReq.accuracy = float32(accuracy64)

		// Adding current timestamp to measure the deadline
		funcReq.timestamp = time.Now()

		// Add the task to pending queue
		queue.addToQueue(funcReq)
	},
	)
	return r
}

func main() {
	// setup_db()
	queue := InitQueue()
	k8s := initKubernetes()
	gpuResources := InitResources(k8s)

	var resourceMutex sync.Mutex

	go queue.schedulingPolicy(k8s, gpuResources, &resourceMutex)
	go monitorPods(k8s, gpuResources, &resourceMutex)

	router := initRouter(queue)
	_ = router.Run("localhost:8083")
}
