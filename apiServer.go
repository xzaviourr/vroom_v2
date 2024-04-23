package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func initServer(reqQueue *ReqQueue, resourceManager *ResourceManager) *gin.Engine {
	r := gin.Default()
	fmt.Println("API server is running")
	r.POST("/run", func(c *gin.Context) {
		var funcReq FuncReq

		if err := c.BindJSON(&funcReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		funcReq.Uid = uuid.New().String()
		funcReq.RegistrationTs = time.Now()
		funcReq.State = "new"
		reqQueue.Enque(&funcReq)

		fmt.Printf("New Request : %s\n", funcReq.Uid)
		response := fmt.Sprintf("Request id: %s", funcReq.Uid)
		c.JSON(http.StatusOK, gin.H{"message": response})
	})

	r.POST("/insert", func(c *gin.Context) {
		var variant Variant

		if err := c.BindJSON(&variant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		variant.Id = uuid.New().String()
		resourceManager.VariantStore.addVariant(&variant)

		c.JSON(http.StatusOK, gin.H{"message": "Variant inserted successfully"})
	})
	return r
}
