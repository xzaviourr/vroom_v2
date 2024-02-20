package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type FuncInfo struct {
	variant_id      string  `json:"variant_id"`
	task_identifier string  `json:"task_identifier"`
	gpu_memory      int     `json:"gpu_memory"`
	gpu_cores       int     `json:"gpu_cores"`
	image           string  `json:"image"`
	latency         float32 `json:"latency"`
	accuracy        float32 `json:"accuracy"`
	batch_size      int     `json:"batch_size"`
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

func initKubernetes() *kubernetes.Clientset {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
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
	queue := InitQueue()
	k8s := initKubernetes()

	go queue.schedulingPolicy(k8s)

	router := initRouter(queue)
	_ = router.Run("localhost:8083")
}
