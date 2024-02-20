package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

type Queue struct {
	items []FuncInfo
}

type k8sClient struct {
	clientset *kubernetes.Clientset
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
		fid := c.Query("id")
		fidInt, _ := strconv.Atoi(fid)
		funcInfo, err := getVariantByID(fidInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "impossible to retrieve"})
			return
		}
		queue.addToQueue(funcInfo)
	},
	)
	return r
}

func systemTimer() int {
	time.Sleep(1 * time.Second)
	i := 0
	i += 1
	return i
}

func main() {
	queue := InitQueue()
	k8s := initKubernetes()

	go systemTimer()
	go queue.schedulingPolicy(k8s)

	router := initRouter(queue)
	_ = router.Run("localhost:8083")
}
