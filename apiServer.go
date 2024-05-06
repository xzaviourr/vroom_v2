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
	fmt.Println("Vroom Server Running. Active End Points : ['/run', '/insert']")

	// Send a new function request at this end point for execution
	r.POST("/run", func(c *gin.Context) {
		var funcReq FuncReq

		if err := c.BindJSON(&funcReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		funcReq.Uid = uuid.New().String() // Generate new unique id for the request
		funcReq.RegistrationTs = time.Now()
		funcReq.State = "new"
		reqQueue.Enque(&funcReq)

		response := fmt.Sprintf("Request registered. Id: %s", funcReq.Uid)
		// Send the request Id back to the user
		c.JSON(http.StatusOK, gin.H{"message": response})
	})

	// Insert a new variant info using this end point
	r.POST("/insert", func(c *gin.Context) {
		var variant Variant

		if err := c.BindJSON(&variant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		variant.Id = uuid.New().String() // Generate new unique id for this variant
		resourceManager.VariantStore.addVariant(&variant)

		response := fmt.Sprintf("Variant stored successfully. Id: %s", variant.Id)
		// Send the variant Id back to the user
		c.JSON(http.StatusOK, gin.H{"message": response})
	})
	return r
}
